//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	mongocontainer "github.com/radiophysiker/microservices-homework/platform/pkg/testcontainers/mongo"
	"github.com/radiophysiker/microservices-homework/platform/pkg/testcontainers/network"
)

// TestEnvironment представляет тестовое окружение с контейнерами
type TestEnvironment struct {
	Network     *network.Network
	Mongo       *mongocontainer.Container
	MongoClient *mongo.Client
	Database    *mongo.Database
	Collection  *mongo.Collection
	AppAddress  string
	MongoURI    string
	appProcess  *exec.Cmd
	envFilePath string // Путь к временному .env файлу
	useCompose  bool   // Использовать Docker Compose вместо testcontainers
	composePath string // Путь к docker-compose.yml
}

// Setup создает и запускает тестовое окружение
func Setup(ctx context.Context) (*TestEnvironment, error) {
	env := &TestEnvironment{}

	// Инициализируем логгер
	if err := logger.Init("info", false); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// Создаем Docker сеть
	testNetwork, err := network.NewNetwork(ctx, "inventory-integration-test")
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}
	env.Network = testNetwork

	// Генерируем уникальное имя контейнера для каждого запуска тестов
	testID := uuid.New().String()[:8]
	mongoContainerName := fmt.Sprintf("%s-%s", TestMongoContainerName, testID)

	// Запускаем MongoDB контейнер
	mongoContainer, err := mongocontainer.NewContainer(
		ctx,
		mongocontainer.WithNetworkName(testNetwork.Name()),
		mongocontainer.WithContainerName(mongoContainerName),
		mongocontainer.WithDatabase(TestMongoDatabase),
		mongocontainer.WithAuth(TestMongoUsername, TestMongoPassword),
		mongocontainer.WithAuthDB(TestMongoAuthDB),
		mongocontainer.WithLogger(&testLogger{}),
	)
	if err != nil {
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}
	env.Mongo = mongoContainer

	// Получаем MongoDB клиент и коллекцию
	env.MongoClient = mongoContainer.Client()
	env.Database = env.MongoClient.Database(TestMongoDatabase)
	env.Collection = env.Database.Collection(TestCollectionName)

	// Получаем MongoDB URI
	mongoCfg := mongoContainer.Config()
	env.MongoURI = fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		mongoCfg.Username,
		mongoCfg.Password,
		mongoCfg.Host,
		mongoCfg.Port,
		mongoCfg.Database,
		mongoCfg.AuthDB,
	)

	// Запускаем приложение напрямую (не в контейнере) для e2e тестов
	// ВАЖНО: Приложение запускается ДО вставки данных, но это нормально,
	// так как данные будут вставлены в SetupSuite перед запуском тестов
	// Устанавливаем переменные окружения для подключения к MongoDB в контейнере
	// Используем localhost для подключения к MongoDB, так как testcontainers маппит порты на localhost
	envVars := []string{
		fmt.Sprintf("GRPC_HOST=0.0.0.0"),
		fmt.Sprintf("GRPC_PORT=%s", TestAppPort),
		fmt.Sprintf("MONGO_HOST=localhost"), // testcontainers маппит порты на localhost
		fmt.Sprintf("MONGO_PORT=%s", mongoCfg.Port),
		fmt.Sprintf("MONGO_DATABASE=%s", TestMongoDatabase),
		fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", TestMongoUsername),
		fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", TestMongoPassword),
		fmt.Sprintf("MONGO_AUTH_DB=%s", TestMongoAuthDB),
		"LOGGER_LEVEL=info",
		"LOGGER_AS_JSON=false",
	}

	// Находим путь к бинарнику приложения
	// Получаем абсолютный путь к директории inventory
	wd, err := os.Getwd()
	if err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Путь к директории inventory (на 2 уровня выше от tests/integration)
	inventoryDir := filepath.Join(wd, "..", "..")
	appPath := filepath.Join(inventoryDir, "main")

	// Если бинарника нет, пытаемся собрать его
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", "main", "./cmd/main.go")
		buildCmd.Dir = inventoryDir
		buildCmd.Env = os.Environ()
		if err := buildCmd.Run(); err != nil {
			_ = mongoContainer.Terminate(ctx)
			_ = testNetwork.Remove(ctx)
			return nil, fmt.Errorf("failed to build app: %w", err)
		}
	}

	// Создаем временный .env файл для приложения, чтобы оно могло загрузить конфигурацию
	envFile := filepath.Join(inventoryDir, ".env")
	envContent := fmt.Sprintf(`GRPC_HOST=%s
GRPC_PORT=%s
MONGO_HOST=%s
MONGO_PORT=%s
MONGO_DATABASE=%s
MONGO_INITDB_ROOT_USERNAME=%s
MONGO_INITDB_ROOT_PASSWORD=%s
MONGO_AUTH_DB=%s
LOGGER_LEVEL=info
LOGGER_AS_JSON=false
`, "0.0.0.0", TestAppPort, "localhost", mongoCfg.Port, TestMongoDatabase, TestMongoUsername, TestMongoPassword, TestMongoAuthDB)

	if err := os.WriteFile(envFile, []byte(envContent), 0o644); err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to create .env file: %w", err)
	}
	env.envFilePath = envFile

	// ВАЖНО: Данные должны быть вставлены ДО запуска приложения
	// Но здесь мы только настраиваем окружение, данные будут вставлены в SetupSuite
	// после вызова Setup

	// Запускаем приложение в фоновом режиме
	appCmd := exec.Command(appPath)
	appCmd.Env = append(os.Environ(), envVars...)
	appCmd.Dir = inventoryDir
	// Перенаправляем вывод в /dev/null, чтобы не засорять вывод тестов
	appCmd.Stdout = nil
	appCmd.Stderr = nil

	if err := appCmd.Start(); err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to start app: %w", err)
	}

	env.appProcess = appCmd
	env.AppAddress = fmt.Sprintf("localhost:%s", TestAppPort)

	// Ждем, чтобы приложение запустилось и было готово принимать соединения
	// Проверяем доступность gRPC сервера
	maxRetries := 60 // Увеличиваем количество попыток
	checkCtx, checkCancel := context.WithTimeout(ctx, 30*time.Second)
	defer checkCancel()

	var conn *grpc.ClientConn
	for i := 0; i < maxRetries; i++ {
		time.Sleep(500 * time.Millisecond) // Увеличиваем задержку между попытками

		// Проверяем, что процесс еще работает
		if env.appProcess.ProcessState != nil && env.appProcess.ProcessState.Exited() {
			_ = mongoContainer.Terminate(ctx)
			_ = testNetwork.Remove(ctx)
			return nil, fmt.Errorf("app process exited unexpectedly")
		}

		// Пытаемся подключиться к gRPC серверу
		var connErr error
		conn, connErr = grpc.NewClient(
			env.AppAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if connErr == nil {
			// Проверяем, что соединение действительно работает
			_ = conn.Close()
			break
		}

		// Проверяем контекст на отмену
		select {
		case <-checkCtx.Done():
			_ = env.appProcess.Process.Kill()
			_ = mongoContainer.Terminate(ctx)
			_ = testNetwork.Remove(ctx)
			return nil, fmt.Errorf("app failed to start: timeout waiting for gRPC server")
		default:
		}

		if i == maxRetries-1 {
			_ = env.appProcess.Process.Kill()
			_ = mongoContainer.Terminate(ctx)
			_ = testNetwork.Remove(ctx)
			return nil, fmt.Errorf("app failed to start: gRPC server not available after %d retries: %v", maxRetries, connErr)
		}
	}

	return env, nil
}

// Teardown останавливает и удаляет тестовое окружение
func (env *TestEnvironment) Teardown(ctx context.Context) error {
	var errs []error

	// Останавливаем приложение
	if env.appProcess != nil {
		if err := env.appProcess.Process.Kill(); err != nil {
			errs = append(errs, fmt.Errorf("failed to kill app process: %w", err))
		}
		// Ждем завершения процесса
		_ = env.appProcess.Wait()
	}

	// Удаляем временный .env файл
	if env.envFilePath != "" {
		if err := os.Remove(env.envFilePath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("failed to remove .env file: %w", err))
		}
	}

	if env.MongoClient != nil {
		if err := env.MongoClient.Disconnect(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to disconnect mongo client: %w", err))
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate mongo container: %w", err))
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove network: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("teardown errors: %v", errs)
	}

	return nil
}

// testLogger реализует интерфейс Logger для testcontainers
type testLogger struct{}

func (l *testLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}

func (l *testLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(ctx, msg, fields...)
}
