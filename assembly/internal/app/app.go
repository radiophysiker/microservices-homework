package app

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/radiophysiker/microservices-homework/assembly/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	"github.com/radiophysiker/microservices-homework/platform/pkg/metrics"
)

type App struct {
	diContainer *diContainer
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

	// Запускаем consumer для OrderPaid
	g.Go(func() error {
		orderConsumerService, err := a.diContainer.OrderConsumerService(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get OrderConsumerService", zap.Error(err))
			return err
		}

		logger.Info(ctx, "Starting OrderPaid consumer")

		if err := orderConsumerService.RunConsumer(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info(ctx, "OrderPaid consumer stopped")
				return nil
			}

			logger.Error(ctx, "OrderPaid consumer error", zap.Error(err))

			return err
		}

		return nil
	})

	// Завершаем по ctx
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		defer cancel()

		logger.Info(shutdownCtx, "Shutting down AssemblyService...")

		return nil
	})

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMetrics,
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

func (a *App) initMetrics(ctx context.Context) error {
	if err := metrics.InitProvider(ctx, config.AppConfig().Metrics); err != nil {
		return err
	}

	closer.AddNamed("Metrics provider", metrics.Shutdown)

	return nil
}
