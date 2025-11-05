package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	apiv1 "github.com/radiophysiker/microservices-homework/inventory/internal/api/inventory/v1"
	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository"
	partRepo "github.com/radiophysiker/microservices-homework/inventory/internal/repository/part"
	"github.com/radiophysiker/microservices-homework/inventory/internal/service"
	partSvc "github.com/radiophysiker/microservices-homework/inventory/internal/service/part"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
)

type diContainer struct {
	mongoClient    *mongo.Client
	collection     *mongo.Collection
	partRepository repository.PartRepository
	partService    service.PartService
	api            *apiv1.API
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) MongoClient(ctx context.Context) (*mongo.Client, error) {
	if d.mongoClient == nil {
		ctxDial, cancelDial := context.WithTimeout(ctx, 15*time.Second)
		defer cancelDial()

		client, err := mongo.Connect(ctxDial, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			return nil, fmt.Errorf("connect mongo: %w", err)
		}

		ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
		defer cancelPing()

		if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
			ctxDisconnect, cancelDisconnect := context.WithTimeout(ctx, 3*time.Second)
			defer cancelDisconnect()

			if disconnectErr := client.Disconnect(ctxDisconnect); disconnectErr != nil {
				return nil, fmt.Errorf("disconnect mongo after ping failure: %w (ping error: %w)", disconnectErr, err)
			}

			return nil, fmt.Errorf("ping mongo: %w", err)
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoClient = client
	}

	return d.mongoClient, nil
}

func (d *diContainer) Collection(ctx context.Context) (*mongo.Collection, error) {
	if d.collection == nil {
		client, err := d.MongoClient(ctx)
		if err != nil {
			return nil, err
		}

		d.collection = client.Database(config.AppConfig().Mongo.DatabaseName()).Collection("parts")
	}

	return d.collection, nil
}

func (d *diContainer) PartRepository(ctx context.Context) (repository.PartRepository, error) {
	if d.partRepository == nil {
		collection, err := d.Collection(ctx)
		if err != nil {
			return nil, err
		}

		d.partRepository = partRepo.NewRepository(collection)
	}

	return d.partRepository, nil
}

func (d *diContainer) PartService(ctx context.Context) (service.PartService, error) {
	if d.partService == nil {
		partRepo, err := d.PartRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.partService = partSvc.NewService(partRepo)
	}

	return d.partService, nil
}

func (d *diContainer) API(ctx context.Context) (*apiv1.API, error) {
	if d.api == nil {
		partService, err := d.PartService(ctx)
		if err != nil {
			return nil, err
		}

		d.api = apiv1.NewAPI(partService)
	}

	return d.api, nil
}
