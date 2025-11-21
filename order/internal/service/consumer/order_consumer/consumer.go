package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/order/internal/converter/kafka/decoder"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type Service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  *decoder.Decoder
	orderRepository        repository.OrderRepository
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder *decoder.Decoder,
	orderRepository repository.OrderRepository,
) *Service {
	return &Service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderRepository:        orderRepository,
	}
}

func (s *Service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderAssembled consumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.ShipAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
