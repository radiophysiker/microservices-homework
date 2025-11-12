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
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
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

// NewGRPCClient создает новый gRPC клиент для тестов
// В новом API соединение устанавливается асинхронно, первый RPC вызов будет ждать установки соединения
func (env *TestEnvironment) NewGRPCClient(ctx context.Context) (pb.InventoryServiceClient, func(), error) {
	conn, err := grpc.NewClient(
		env.AppAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	client := pb.NewInventoryServiceClient(conn)
	cleanup := func() {
		_ = conn.Close()
	}

	return client, cleanup, nil
}

// Setup создает и запускает тестовое окружение
func Setup(ctx context.Context) (*TestEnvironment, error) {
	env := &TestEnvironment{}

	if err := logger.Init("info", false); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	testNetwork, err := network.NewNetwork(ctx, "inventory-integration-test")
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}
	env.Network = testNetwork

	testID := uuid.New().String()[:8]
	mongoContainerName := fmt.Sprintf("%s-%s", TestMongoContainerName, testID)

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

	env.MongoClient = mongoContainer.Client()
	env.Database = env.MongoClient.Database(TestMongoDatabase)
	env.Collection = env.Database.Collection(TestCollectionName)

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

	wd, err := os.Getwd()
	if err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	inventoryDir := filepath.Join(wd, "..", "..")
	appPath := filepath.Join(inventoryDir, "main")

	buildCmd := exec.Command("go", "build", "-o", "main", "./cmd/main.go")
	buildCmd.Dir = inventoryDir
	buildCmd.Env = os.Environ()
	if err := buildCmd.Run(); err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to build app: %w", err)
	}

	appCmd := exec.Command(appPath)
	appCmd.Env = append(os.Environ(), envVars...)
	appCmd.Dir = inventoryDir
	appCmd.Stdout = nil
	appCmd.Stderr = nil

	if err := appCmd.Start(); err != nil {
		_ = mongoContainer.Terminate(ctx)
		_ = testNetwork.Remove(ctx)
		return nil, fmt.Errorf("failed to start app: %w", err)
	}

	env.appProcess = appCmd
	env.AppAddress = fmt.Sprintf("localhost:%s", TestAppPort)

	maxRetries := 60
	checkCtx, checkCancel := context.WithTimeout(ctx, 30*time.Second)
	defer checkCancel()

	var conn *grpc.ClientConn
	for i := 0; i < maxRetries; i++ {
		time.Sleep(500 * time.Millisecond)

		if env.appProcess.ProcessState != nil && env.appProcess.ProcessState.Exited() {
			_ = mongoContainer.Terminate(ctx)
			_ = testNetwork.Remove(ctx)
			return nil, fmt.Errorf("app process exited unexpectedly")
		}

		var connErr error
		conn, connErr = grpc.NewClient(
			env.AppAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if connErr == nil {
			_ = conn.Close()
			break
		}

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

	if env.appProcess != nil {
		if err := env.appProcess.Process.Kill(); err != nil {
			errs = append(errs, fmt.Errorf("failed to kill app process: %w", err))
		}
		_ = env.appProcess.Wait()
	}

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
