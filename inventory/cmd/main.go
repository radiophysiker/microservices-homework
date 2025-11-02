package main

import (
	"context"
	"log"

	"github.com/radiophysiker/microservices-homework/inventory/internal/app"
	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
)

const configPath = ".env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
