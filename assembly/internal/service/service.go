package service

import (
	"context"

	"github.com/radiophysiker/microservices-homework/assembly/internal/model"
)

// OrderConsumerService представляет интерфейс для consumer'а событий OrderPaid
type OrderConsumerService interface {
	// RunConsumer запускает consumer для обработки событий OrderPaid
	RunConsumer(ctx context.Context) error
}

// ShipAssembledProducerService представляет интерфейс для producer'а событий ShipAssembled
type ShipAssembledProducerService interface {
	// ProduceShipAssembled отправляет событие ShipAssembled в Kafka
	ProduceShipAssembled(ctx context.Context, shipAssembled model.ShipAssembled) error
}
