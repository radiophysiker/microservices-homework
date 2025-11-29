package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/radiophysiker/microservices-homework/assembly/internal/config"
	kafkaConverter "github.com/radiophysiker/microservices-homework/assembly/internal/converter/kafka"
	"github.com/radiophysiker/microservices-homework/assembly/internal/converter/kafka/decoder"
	svc "github.com/radiophysiker/microservices-homework/assembly/internal/service"
	orderConsumerSvc "github.com/radiophysiker/microservices-homework/assembly/internal/service/consumer/order_consumer"
	orderProducerSvc "github.com/radiophysiker/microservices-homework/assembly/internal/service/producer/order_producer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	kafkaConsumer "github.com/radiophysiker/microservices-homework/platform/pkg/kafka/consumer"
	kafkaProducer "github.com/radiophysiker/microservices-homework/platform/pkg/kafka/producer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer      kafka.Consumer

	shipAssembledSyncProducer sarama.SyncProducer
	shipAssembledProducer     kafka.Producer

	orderPaidDecoder kafkaConverter.OrderPaidDecoder

	orderConsumerService         svc.OrderConsumerService
	shipAssembledProducerService svc.ShipAssembledProducerService
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerGroup(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.orderPaidConsumerGroup == nil {
		cfg := config.AppConfig()
		consumerCfg := cfg.OrderPaidConsumer

		group, err := sarama.NewConsumerGroup(
			cfg.Kafka.Brokers(),
			consumerCfg.GroupID(),
			consumerCfg.Config(),
		)
		if err != nil {
			return nil, fmt.Errorf("create consumer group: %w", err)
		}

		closer.AddNamed("OrderPaid consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.orderPaidConsumerGroup = group
	}

	return d.orderPaidConsumerGroup, nil
}

func (d *diContainer) ShipAssembledSyncProducer(ctx context.Context) (sarama.SyncProducer, error) {
	if d.shipAssembledSyncProducer == nil {
		cfg := config.AppConfig()
		producerCfg := cfg.OrderAssembledProducer

		producer, err := sarama.NewSyncProducer(
			cfg.Kafka.Brokers(),
			producerCfg.Config(),
		)
		if err != nil {
			return nil, fmt.Errorf("create sync producer: %w", err)
		}

		closer.AddNamed("ShipAssembled sync producer", func(ctx context.Context) error {
			return producer.Close()
		})

		d.shipAssembledSyncProducer = producer
	}

	return d.shipAssembledSyncProducer, nil
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) (kafka.Consumer, error) {
	if d.orderPaidConsumer == nil {
		group, err := d.OrderPaidConsumerGroup(ctx)
		if err != nil {
			return nil, err
		}

		cfg := config.AppConfig()
		topics := []string{cfg.OrderPaidConsumer.Topic()}

		d.orderPaidConsumer = kafkaConsumer.NewConsumer(
			group,
			topics,
			logger.Logger(),
		)
	}

	return d.orderPaidConsumer, nil
}

func (d *diContainer) ShipAssembledProducer(ctx context.Context) (kafka.Producer, error) {
	if d.shipAssembledProducer == nil {
		syncProducer, err := d.ShipAssembledSyncProducer(ctx)
		if err != nil {
			return nil, err
		}

		cfg := config.AppConfig()
		topic := cfg.OrderAssembledProducer.Topic()

		d.shipAssembledProducer = kafkaProducer.NewProducer(
			syncProducer,
			topic,
			logger.Logger(),
		)
	}

	return d.shipAssembledProducer, nil
}

func (d *diContainer) OrderPaidDecoder(_ context.Context) (kafkaConverter.OrderPaidDecoder, error) {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder, nil
}

func (d *diContainer) ShipAssembledProducerService(ctx context.Context) (svc.ShipAssembledProducerService, error) {
	if d.shipAssembledProducerService == nil {
		producer, err := d.ShipAssembledProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.shipAssembledProducerService = orderProducerSvc.NewService(producer)
	}

	return d.shipAssembledProducerService, nil
}

func (d *diContainer) OrderConsumerService(ctx context.Context) (svc.OrderConsumerService, error) {
	if d.orderConsumerService == nil {
		consumer, err := d.OrderPaidConsumer(ctx)
		if err != nil {
			return nil, err
		}

		decoder, err := d.OrderPaidDecoder(ctx)
		if err != nil {
			return nil, err
		}

		producerService, err := d.ShipAssembledProducerService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderConsumerService = orderConsumerSvc.NewService(
			ctx,
			consumer,
			decoder,
			producerService,
		)
	}

	return d.orderConsumerService, nil
}
