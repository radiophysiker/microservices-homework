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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	grpcMiddleware "github.com/radiophysiker/microservices-homework/platform/pkg/middleware/grpc"
	httpMiddleware "github.com/radiophysiker/microservices-homework/platform/pkg/middleware/http"
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

	g.Go(func() error {
		orderConsumerService, err := a.diContainer.OrderConsumerService(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get OrderConsumerService", zap.Error(err))
			return err
		}

		logger.Info(ctx, "Starting OrderAssembled consumer")

		if err := orderConsumerService.RunConsumer(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info(ctx, "OrderAssembled consumer stopped")
				return nil
			}

			logger.Error(ctx, "OrderAssembled consumer error", zap.Error(err))

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

	authInterceptor, err := a.diContainer.AuthInterceptor(ctx)
	if err != nil {
		return fmt.Errorf("failed to create auth interceptor: %w", err)
	}

	a.grpcServer = grpc.NewServer(
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

	orderpb.RegisterOrderServiceServer(a.grpcServer, api)

	return nil
}

func (a *App) initHTTPGateway(ctx context.Context) error {
	gatewayCtx, gatewayCancel := context.WithCancel(ctx)
	a.gatewayCancel = gatewayCancel

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.MD{}

			if sessionUUID, ok := grpcMiddleware.GetSessionUUIDFromContext(ctx); ok && sessionUUID != "" {
				md.Set(grpcMiddleware.SessionUUIDMetadataKey, sessionUUID)
			}

			return md
		}),
	)

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

	var authMiddleware *httpMiddleware.AuthMiddleware

	authMiddleware, err = a.diContainer.AuthMiddleware(ctx)
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %w", err)
	}

	handler := authMiddleware.Handle(mux)

	httpAddr := config.AppConfig().OrderHTTP.Address()
	a.httpServer = &http.Server{
		Addr:              httpAddr,
		Handler:           handler,
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
