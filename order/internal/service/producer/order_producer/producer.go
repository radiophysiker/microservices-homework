package order_producer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/order/internal/converter/kafka/encoder"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type Service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *Service {
	return &Service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *Service) ProduceOrderPaid(ctx context.Context, orderPaid model.OrderPaid) error {
	value, err := encoder.EncodeOrderPaid(orderPaid)
	if err != nil {
		logger.Error(ctx, "Failed to encode OrderPaid event",
			zap.Error(err),
			zap.String("order_uuid", orderPaid.OrderUUID.String()),
			zap.String("event_uuid", orderPaid.EventUUID.String()),
		)

		return fmt.Errorf("failed to encode OrderPaid: %w", err)
	}

	key := []byte(orderPaid.OrderUUID.String())

	if err := s.orderPaidProducer.Send(ctx, key, value); err != nil {
		logger.Error(ctx, "Failed to send OrderPaid event",
			zap.Error(err),
			zap.String("order_uuid", orderPaid.OrderUUID.String()),
			zap.String("event_uuid", orderPaid.EventUUID.String()),
		)

		return fmt.Errorf("failed to send OrderPaid event: %w", err)
	}

	logger.Info(ctx, "OrderPaid event sent",
		zap.String("order_uuid", orderPaid.OrderUUID.String()),
		zap.String("event_uuid", orderPaid.EventUUID.String()),
	)

	return nil
}
