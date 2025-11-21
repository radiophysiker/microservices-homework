package order_consumer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

func (s *Service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled event",
			zap.Error(err),
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
		)

		return err
	}

	logger.Info(ctx, "ShipAssembled message received",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("user_uuid", event.UserUUID.String()),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	order, err := s.orderRepository.GetOrder(ctx, event.OrderUUID.String())
	if err != nil {
		logger.Error(ctx, "Failed to get order",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return fmt.Errorf("failed to get order: %w", err)
	}

	order.Status = model.StatusAssembled

	updated, err := s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return fmt.Errorf("failed to update order status: %w", err)
	}

	logger.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", updated.OrderUUID.String()),
		zap.String("event_uuid", event.EventUUID.String()),
	)

	return nil
}
