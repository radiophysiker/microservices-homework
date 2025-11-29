package app

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/radiophysiker/microservices-homework/payment/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	"github.com/radiophysiker/microservices-homework/platform/pkg/tracing"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
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
	logger.Info(ctx, "PaymentService gRPC server listening", zap.String("address", a.listener.Addr().String()))

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
		a.initTracing,
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

func (a *App) initTracing(ctx context.Context) error {
	if err := tracing.InitTracer(ctx, config.AppConfig().Tracing); err != nil {
		return err
	}

	closer.AddNamed("Tracer", tracing.ShutdownTracer)

	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			tracing.UnaryServerInterceptor(config.AppConfig().Tracing.ServiceName()),
		),
	)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	pb.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.API())

	return nil
}
