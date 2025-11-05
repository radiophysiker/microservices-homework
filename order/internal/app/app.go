package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	"github.com/radiophysiker/microservices-homework/platform/pkg/migrator"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

type App struct {
	diContainer   *diContainer
	grpcServer    *grpc.Server
	httpServer    *http.Server
	grpcListener  net.Listener
	gatewayCancel context.CancelFunc
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
	parentCtx := ctx
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info(ctx, "OrderService gRPC listen", zap.String("addr", a.grpcListener.Addr().String()))

		if err := a.grpcServer.Serve(a.grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error(ctx, "gRPC serve failed", zap.Error(err))
			return err
		}

		return nil
	})

	g.Go(func() error {
		logger.Info(ctx, "OrderService HTTP Gateway listen", zap.String("addr", a.httpServer.Addr))

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "HTTP serve failed", zap.Error(err))
			return err
		}

		return nil
	})

	// Завершаем по ctx
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		defer cancel()
		a.gatewayCancel()
		a.grpcServer.GracefulStop()

		return a.httpServer.Shutdown(shutdownCtx)
	})

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMigrations,
		a.initGRPCServer,
		a.initHTTPGateway,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	pool, err := a.diContainer.Pool(ctx)
	if err != nil {
		return err
	}

	return migrator.Run(ctx, pool, config.AppConfig().Migrations.Directory())
}

func (a *App) initGRPCServer(ctx context.Context) error {
	grpcAddr := config.AppConfig().OrderGRPC.Address()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	a.grpcListener = lis

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := lis.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.grpcServer = grpc.NewServer()

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

	orderpb.RegisterOrderServiceServer(a.grpcServer, api)

	return nil
}

func (a *App) initHTTPGateway(ctx context.Context) error {
	gatewayCtx, gatewayCancel := context.WithCancel(ctx)
	a.gatewayCancel = gatewayCancel

	mux := runtime.NewServeMux()

	err := orderpb.RegisterOrderServiceHandlerFromEndpoint(
		gatewayCtx,
		mux,
		config.AppConfig().OrderGRPC.Address(),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	httpAddr := config.AppConfig().OrderHTTP.Address()
	a.httpServer = &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		a.gatewayCancel()

		shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		return a.httpServer.Shutdown(shutdownCtx)
	})

	return nil
}
