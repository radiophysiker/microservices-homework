package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event",
			zap.Error(err),
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
		)

		return err
	}

	logger.Info(ctx, "OrderPaid message received",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("user_uuid", event.UserUUID.String()),
		zap.String("transaction_uuid", event.TransactionUUID.String()),
		zap.Int("payment_method", int(event.PaymentMethod)),
	)

	if err := s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send OrderPaid notification",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return err
	}

	return nil
}
