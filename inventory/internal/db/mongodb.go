package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
)

// Connect - Соединяется с MongoDB и возвращает клиент и коллекцию
// ctx - контекст
// cfg - конфигурация
// Возвращает клиент, коллекцию и ошибку
func Connect(ctx context.Context) (*mongo.Client, *mongo.Collection, error) {
	ctxDial, cancelDial := context.WithTimeout(ctx, 15*time.Second)
	defer cancelDial()

	client, err := mongo.Connect(ctxDial, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		return nil, nil, fmt.Errorf("connect mongo: %w", err)
	}

	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		ctxDisconnect, cancelDisconnect := context.WithTimeout(ctx, 3*time.Second)
		defer cancelDisconnect()

		if disconnectErr := client.Disconnect(ctxDisconnect); disconnectErr != nil {
			log.Printf("failed to disconnect MongoDB client: %v", disconnectErr)
		}

		return nil, nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database(config.AppConfig().Mongo.DatabaseName())
	collection := db.Collection("parts")

	return client, collection, nil
}
