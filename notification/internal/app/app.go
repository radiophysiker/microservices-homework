package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/radiophysiker/microservices-homework/notification/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
		logger.Info(ctx, "NotificationService HTTP server listen", zap.String("addr", a.httpServer.Addr))

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "HTTP serve failed", zap.Error(err))
			return err
		}

		return nil
	})

	g.Go(func() error {
		orderPaidConsumerService, err := a.diContainer.OrderPaidConsumerService(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get OrderPaidConsumerService", zap.Error(err))
			return err
		}

		logger.Info(ctx, "Starting OrderPaid consumer")

		if err := orderPaidConsumerService.RunConsumer(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info(ctx, "OrderPaid consumer stopped")
				return nil
			}

			logger.Error(ctx, "OrderPaid consumer error", zap.Error(err))

			return err
		}

		return nil
	})

	g.Go(func() error {
		orderAssembledConsumerService, err := a.diContainer.OrderAssembledConsumerService(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get OrderAssembledConsumerService", zap.Error(err))
			return err
		}

		logger.Info(ctx, "Starting OrderAssembled consumer")

		if err := orderAssembledConsumerService.RunConsumer(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info(ctx, "OrderAssembled consumer stopped")
				return nil
			}

			logger.Error(ctx, "OrderAssembled consumer error", zap.Error(err))

			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		defer cancel()

		logger.Info(shutdownCtx, "Shutting down NotificationService...")

		return a.httpServer.Shutdown(shutdownCtx)
	})

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) error {
	api, err := a.diContainer.API(ctx)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	httpAddr := config.AppConfig().HTTP.Address()
	a.httpServer = &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		return a.httpServer.Shutdown(shutdownCtx)
	})

	return nil
}
