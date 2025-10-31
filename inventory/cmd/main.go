package main

import (
	"context"
	"log"

	"github.com/radiophysiker/microservices-homework/inventory/internal/app"
	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
