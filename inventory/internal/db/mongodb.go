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

// Соединяется с MongoDB и возвращает клиент и коллекцию
// ctx - контекст
// cfg - конфигурация
// Возвращает клиент, коллекцию и ошибку
func Connect(ctx context.Context, cfg config.Config) (*mongo.Client, *mongo.Collection, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.DBHost, cfg.DBPort)

	clientOpts := options.Client().
		ApplyURI(uri).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(20).
		SetMinPoolSize(2)

	if cfg.DBUser != "" && cfg.DBPass != "" {
		clientOpts.SetAuth(options.Credential{
			AuthSource: cfg.DBAuthDB,
			Username:   cfg.DBUser,
			Password:   cfg.DBPass,
		})
	}

	ctxDial, cancelDial := context.WithTimeout(ctx, 15*time.Second)
	defer cancelDial()

	client, err := mongo.Connect(ctxDial, clientOpts)
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

	db := client.Database(cfg.DBName)
	collection := db.Collection("parts")

	return client, collection, nil
}
