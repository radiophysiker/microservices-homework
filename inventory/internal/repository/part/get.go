package part

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
)

// GetPart возвращает деталь по UUID
func (r *Repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	filter := bson.M{"uuid": uuid}

	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, filter).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("part not found")
		}

		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return converter.ToServicePart(&repoPart), nil
}
