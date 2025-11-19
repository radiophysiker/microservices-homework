package order_consumer

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/assembly/internal/model"
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

	delaySeconds, err := getRandomDelay1to10()
	if err != nil {
		logger.Error(ctx, "Failed to get random delay",
			zap.Error(err),
		)

		return err
	}

	delay := time.Duration(delaySeconds) * time.Second

	logger.Info(ctx, "Starting ship assembly",
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.Duration("delay", delay),
		zap.Int("delay_seconds", delaySeconds),
	)

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Ship assembly cancelled",
			zap.String("order_uuid", event.OrderUUID.String()),
			zap.Error(ctx.Err()),
		)

		return ctx.Err()
	case <-time.After(delay):
	}

	// Создаем событие ShipAssembled
	shipAssembled := model.ShipAssembled{
		EventUUID:    uuid.New(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(delaySeconds),
	}

	// Публикуем событие ShipAssembled
	if err := s.shipAssembledProducer.ProduceShipAssembled(ctx, shipAssembled); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled event",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return err
	}

	logger.Info(ctx, "Ship assembly completed and ShipAssembled event published",
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("ship_assembled_event_uuid", shipAssembled.EventUUID.String()),
	)

	return nil
}

func getRandomDelay1to10() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + 1, nil
}
