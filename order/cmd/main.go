package main

import (
	"context"
	"log"

	"github.com/radiophysiker/microservices-homework/order/internal/app"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
)

const configPath = ".env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("app: %v", err)
	}
}
