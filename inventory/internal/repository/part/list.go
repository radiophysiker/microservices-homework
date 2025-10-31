package part

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
)

// ListParts возвращает список деталей с возможностью фильтрации
func (r *Repository) ListParts(ctx context.Context, filter *model.Filter) ([]*model.Part, error) {
	mongoFilter := r.buildMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}()

	var repoParts []*repoModel.Part
	if err := cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	return converter.ToServiceParts(repoParts), nil
}

func (r *Repository) buildMongoFilter(filter *model.Filter) bson.M {
	if filter == nil {
		return bson.M{}
	}

	mongoFilter := bson.M{}

	if len(filter.UUIDs) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		repoCategories := make([]repoModel.Category, len(filter.Categories))
		for i, cat := range filter.Categories {
			repoCategories[i] = converter.ToRepoCategory(cat)
		}

		mongoFilter["category"] = bson.M{"$in": repoCategories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
