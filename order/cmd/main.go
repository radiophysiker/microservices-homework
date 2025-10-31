package main

import (
	"context"
	"log"

	"github.com/radiophysiker/microservices-homework/order/internal/app"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := app.Run(context.Background(), cfg); err != nil {
		log.Fatalf("app: %v", err)
	}
}
