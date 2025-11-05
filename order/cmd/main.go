package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/order/internal/app"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

const configPath = "./deploy/compose/order/.env"

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	appInstance, err := app.New(appCtx)
	if err != nil {
		logger.SetNopLogger()
		logger.Fatal(appCtx, "failed to create app", zap.Error(err))
	}

	if err := appInstance.Run(appCtx); err != nil {
		logger.Fatal(appCtx, "failed to run app", zap.Error(err))
	}
}

func gracefulShutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		logger.Fatal(context.Background(), "failed to shutdown gracefully", zap.Error(err))
	}
}
