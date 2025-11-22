package app

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/radiophysiker/microservices-homework/iam/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	"github.com/radiophysiker/microservices-homework/platform/pkg/migrator"
	authpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/auth/v1"
	userpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/user/v1"
)

// App представляет основное приложение IAM сервиса.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New создает новый экземпляр App и инициализирует все зависимости.
func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает gRPC сервер и обрабатывает входящие запросы.
// Блокирует выполнение до остановки сервера или ошибки.
func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "IAMService gRPC server listening", zap.String("address", a.listener.Addr().String()))

	if err := a.grpcServer.Serve(a.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Error(ctx, "gRPC serve failed", zap.Error(err))
		return err
	}

	return nil
}

// initDeps инициализирует все зависимости приложения в правильном порядке.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMigrations,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initDI инициализирует DI контейнер.
func (a *App) initDI(_ context.Context) error {
	a.diContainer = newDiContainer()
	return nil
}

// initLogger инициализирует логгер на основе конфигурации.
func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

// initCloser настраивает closer для graceful shutdown.
func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

// initMigrations применяет миграции базы данных при старте приложения.
func (a *App) initMigrations(ctx context.Context) error {
	pool, err := a.diContainer.Pool(ctx)
	if err != nil {
		return err
	}

	return migrator.Run(ctx, pool, config.AppConfig().Migrations.Directory())
}

// initListener создает TCP listener для gRPC сервера.
func (a *App) initListener(_ context.Context) error {
	addr := config.AppConfig().IAMGRPC.Address()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	a.listener = lis

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := lis.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	return nil
}

// initGRPCServer инициализирует gRPC сервер, регистрирует сервисы и настраивает health checks.
func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)

	authAPI, err := a.diContainer.AuthAPI(ctx)
	if err != nil {
		return err
	}

	userAPI, err := a.diContainer.UserAPI(ctx)
	if err != nil {
		return err
	}

	authpb.RegisterAuthServiceServer(a.grpcServer, authAPI)
	userpb.RegisterUserServiceServer(a.grpcServer, userAPI)

	logger.Info(ctx, "gRPC server initialized")

	return nil
}
