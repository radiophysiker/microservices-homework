package order_consumer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
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
	assemblyDuration      metric.Float64Histogram
}

func NewService(
	ctx context.Context,
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	shipAssembledProducer svc.ShipAssembledProducerService,
) svc.OrderConsumerService {
	meter := otel.Meter("assembly-service")

	assemblyDuration, err := meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Time taken to assemble a ship"),
		metric.WithUnit("s"),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create assembly_duration_seconds histogram", zap.Error(err))
	}

	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		orderPaidDecoder:      orderPaidDecoder,
		shipAssembledProducer: shipAssembledProducer,
		assemblyDuration:      assemblyDuration,
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
