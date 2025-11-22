package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/part"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	clearFlag := flag.Bool("clear", false, "Clear existing data before seeding")

	flag.Parse()

	if err := config.Load(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(config.AppConfig().Logger.Level(), config.AppConfig().Logger.AsJson()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to MongoDB", zap.Error(err))
	}

	disconnectBaseCtx := context.WithoutCancel(ctx)
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(disconnectBaseCtx, 5*time.Second)
		defer disconnectCancel()

		if err := client.Disconnect(disconnectCtx); err != nil {
			logger.Error(disconnectCtx, "Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Проверка соединения
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		logger.Fatal(ctx, "Failed to ping MongoDB", zap.Error(err))
	}

	collection := client.Database(config.AppConfig().Mongo.DatabaseName()).Collection("parts")

	// Очистка данных, если указан флаг
	if *clearFlag {
		logger.Info(ctx, "Clearing existing data")

		if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
			logger.Fatal(ctx, "Failed to clear existing data", zap.Error(err))
		}

		logger.Info(ctx, "Existing data cleared")
	}

	// Получение тестовых данных
	testParts := part.GetTestParts()
	logger.Info(ctx, "Seeding test data", zap.Int("count", len(testParts)))

	// Вставка данных
	repoParts := make([]interface{}, 0, len(testParts))

	for _, part := range testParts {
		repoPart := converter.ToRepoPart(part)
		if repoPart.CreatedAt.IsZero() {
			repoPart.CreatedAt = time.Now()
		}

		if repoPart.UpdatedAt.IsZero() {
			repoPart.UpdatedAt = time.Now()
		}

		bsonData, err := bson.Marshal(repoPart)
		if err != nil {
			logger.Fatal(ctx, "Failed to marshal part to BSON", zap.Error(err), zap.String("uuid", part.UUID))
		}

		var bsonDoc bson.M
		if err := bson.Unmarshal(bsonData, &bsonDoc); err != nil {
			logger.Fatal(ctx, "Failed to unmarshal BSON", zap.Error(err), zap.String("uuid", part.UUID))
		}

		repoParts = append(repoParts, bsonDoc)
	}

	result, err := collection.InsertMany(ctx, repoParts)
	if err != nil {
		logger.Fatal(ctx, "Failed to insert test data", zap.Error(err))
	}

	logger.Info(ctx, "Test data seeded successfully",
		zap.Int("inserted", len(result.InsertedIDs)),
		zap.Int("expected", len(testParts)),
	)

	// Проверка вставленных данных
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Fatal(ctx, "Failed to count documents", zap.Error(err))
	}

	logger.Info(ctx, "Database seeding completed",
		zap.Int64("total_documents", count),
	)
}
