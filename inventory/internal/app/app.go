package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "InventoryService gRPC server listening", zap.String("address", a.listener.Addr().String()))

	if err := a.grpcServer.Serve(a.listener); err != nil {
		logger.Fatal(ctx, "failed to serve gRPC", zap.Error(err))
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
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

func (a *App) initDI(_ context.Context) error {
	a.diContainer = newDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	if err := logger.Init(ctx, config.AppConfig().Logger); err != nil {
		return err
	}

	closer.AddNamed("OTLP logger exporter", logger.Shutdown)

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	addr := config.AppConfig().InventoryGRPC.Address()

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

func (a *App) initGRPCServer(ctx context.Context) error {
	authInterceptor, err := a.diContainer.AuthInterceptor(ctx)
	if err != nil {
		return fmt.Errorf("failed to create auth interceptor: %w", err)
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)

	api, err := a.diContainer.API(ctx)
	if err != nil {
		return err
	}

	pb.RegisterInventoryServiceServer(a.grpcServer, api)

	return nil
}
