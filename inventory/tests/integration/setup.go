//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
)

// SetupTestData вставляет тестовые данные в MongoDB
func SetupTestData(ctx context.Context, collection *mongo.Collection, parts []*model.Part) error {
	if len(parts) == 0 {
		return nil
	}

	repoParts := make([]interface{}, 0, len(parts))
	for _, part := range parts {
		repoPart := converter.ToRepoPart(part)
		if repoPart.CreatedAt.IsZero() {
			repoPart.CreatedAt = time.Now()
		}
		if repoPart.UpdatedAt.IsZero() {
			repoPart.UpdatedAt = time.Now()
		}

		bsonData, err := bson.Marshal(repoPart)
		if err != nil {
			return fmt.Errorf("failed to marshal part to BSON: %w", err)
		}

		var bsonDoc bson.M
		if err := bson.Unmarshal(bsonData, &bsonDoc); err != nil {
			return fmt.Errorf("failed to unmarshal BSON: %w", err)
		}

		repoParts = append(repoParts, bsonDoc)
	}

	result, err := collection.InsertMany(ctx, repoParts)
	if err != nil {
		return fmt.Errorf("failed to insert test data: %w", err)
	}

	if len(result.InsertedIDs) != len(parts) {
		return fmt.Errorf("expected to insert %d parts, but inserted %d", len(parts), len(result.InsertedIDs))
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}
	if count != int64(len(parts)) {
		return fmt.Errorf("expected %d documents in collection, but found %d", len(parts), count)
	}

	for _, part := range parts {
		var result bson.M
		err := collection.FindOne(ctx, bson.M{"uuid": part.UUID}).Decode(&result)
		if err != nil {
			return fmt.Errorf("failed to find inserted part with uuid %s: %w", part.UUID, err)
		}
	}

	return nil
}

// CleanupTestData удаляет все тестовые данные из коллекции
func CleanupTestData(ctx context.Context, collection *mongo.Collection) error {
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to cleanup test data: %w", err)
	}

	return nil
}

// GetTestParts возвращает тестовые данные деталей для использования в e2e тестах
func GetTestParts() []*model.Part {
	return []*model.Part{
		{
			UUID:        "123e4567-e89b-12d3-a456-426614174000",
			Name:        "Main Engine V8",
			Description: "High-performance rocket engine",
			Price:       50000,
			Category:    model.CategoryEngine,
			Manufacturer: &model.Manufacturer{
				Name:    "SpaceTech",
				Country: "USA",
				Website: "https://spacetech.com",
			},
			Tags: []string{"engine", "propulsion", "v8"},
		},
		{
			UUID:        "223e4567-e89b-12d3-a456-426614174001",
			Name:        "Fuel Tank",
			Description: "Large capacity fuel storage",
			Price:       15000,
			Category:    model.CategoryFuel,
			Manufacturer: &model.Manufacturer{
				Name:    "FuelCorp",
				Country: "Germany",
				Website: "https://fuelcorp.com",
			},
			Tags: []string{"fuel", "storage", "tank"},
		},
		{
			UUID:        "323e4567-e89b-12d3-a456-426614174002",
			Name:        "Wing Assembly",
			Description: "Aerodynamic wing structure",
			Price:       25000,
			Category:    model.CategoryWing,
			Manufacturer: &model.Manufacturer{
				Name:    "AeroParts",
				Country: "France",
				Website: "https://aeroparts.com",
			},
			Tags: []string{"wing", "structure", "aerodynamics"},
		},
		{
			UUID:        "423e4567-e89b-12d3-a456-426614174003",
			Name:        "Cockpit Module",
			Description: "Pilot control center",
			Price:       35000,
			Category:    model.CategoryPorthole,
			Manufacturer: &model.Manufacturer{
				Name:    "ControlTech",
				Country: "Japan",
				Website: "https://controltech.com",
			},
			Tags: []string{"cockpit", "control", "pilot"},
		},
	}
}
