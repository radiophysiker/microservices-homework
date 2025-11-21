package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/radiophysiker/microservices-homework/assembly/internal/converter/kafka"
	svc "github.com/radiophysiker/microservices-homework/assembly/internal/service"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer     kafka.Consumer
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	shipAssembledProducer svc.ShipAssembledProducerService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	shipAssembledProducer svc.ShipAssembledProducerService,
) svc.OrderConsumerService {
	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		orderPaidDecoder:      orderPaidDecoder,
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ufoRecordedConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
