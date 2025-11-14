package order_producer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/assembly/internal/converter/kafka/encoder"
	"github.com/radiophysiker/microservices-homework/assembly/internal/model"
	svc "github.com/radiophysiker/microservices-homework/assembly/internal/service"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type service struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) svc.ShipAssembledProducerService {
	return &service{shipAssembledProducer: shipAssembledProducer}
}

// ProduceShipAssembled отправляет событие ShipAssembled в Kafka
func (s *service) ProduceShipAssembled(ctx context.Context, shipAssembled model.ShipAssembled) error {
	value, err := encoder.EncodeShipAssembled(shipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to encode ShipAssembled event",
			zap.Error(err),
			zap.String("order_uuid", shipAssembled.OrderUUID.String()),
			zap.String("event_uuid", shipAssembled.EventUUID.String()),
		)

		return fmt.Errorf("failed to encode ShipAssembled: %w", err)
	}

	key := []byte(shipAssembled.OrderUUID.String())

	if err := s.shipAssembledProducer.Send(ctx, key, value); err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled event",
			zap.Error(err),
			zap.String("order_uuid", shipAssembled.OrderUUID.String()),
			zap.String("event_uuid", shipAssembled.EventUUID.String()),
		)

		return fmt.Errorf("failed to send ShipAssembled event: %w", err)
	}

	logger.Info(ctx, "ShipAssembled event sent",
		zap.String("order_uuid", shipAssembled.OrderUUID.String()),
		zap.String("event_uuid", shipAssembled.EventUUID.String()),
		zap.Int64("build_time_sec", shipAssembled.BuildTimeSec),
	)

	return nil
}
