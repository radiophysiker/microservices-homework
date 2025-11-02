package main

import (
	"context"
	"log"

	"github.com/radiophysiker/microservices-homework/payment/internal/app"
	"github.com/radiophysiker/microservices-homework/payment/internal/config"
)

const configPath = ".env"

func main() {
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
