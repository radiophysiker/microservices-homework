package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
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

	if err := s.telegramService.SendShipAssembledNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled notification",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return err
	}

	return nil
}
