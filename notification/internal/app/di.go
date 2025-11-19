package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	v1 "github.com/radiophysiker/microservices-homework/notification/internal/api/telegram/v1"
	"github.com/radiophysiker/microservices-homework/notification/internal/client/http/telegram"
	"github.com/radiophysiker/microservices-homework/notification/internal/config"
	kafkaConverter "github.com/radiophysiker/microservices-homework/notification/internal/converter/kafka"
	"github.com/radiophysiker/microservices-homework/notification/internal/converter/kafka/decoder"
	svc "github.com/radiophysiker/microservices-homework/notification/internal/service"
	orderAssembledConsumerSvc "github.com/radiophysiker/microservices-homework/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumerSvc "github.com/radiophysiker/microservices-homework/notification/internal/service/consumer/order_paid_consumer"
	telegramSvc "github.com/radiophysiker/microservices-homework/notification/internal/service/telegram"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	kafkaConsumer "github.com/radiophysiker/microservices-homework/platform/pkg/kafka/consumer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumerGroup      sarama.ConsumerGroup
	orderPaidConsumer           kafka.Consumer
	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      kafka.Consumer

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder

	telegramClient  *telegram.Client
	telegramService svc.TelegramService

	orderPaidConsumerService      svc.OrderPaidConsumerService
	orderAssembledConsumerService svc.OrderAssembledConsumerService

	api *v1.API
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

func (d *diContainer) OrderAssembledConsumerGroup(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.orderAssembledConsumerGroup == nil {
		cfg := config.AppConfig()
		consumerCfg := cfg.OrderAssembledConsumer

		group, err := sarama.NewConsumerGroup(
			cfg.Kafka.Brokers(),
			consumerCfg.GroupID(),
			consumerCfg.Config(),
		)
		if err != nil {
			return nil, fmt.Errorf("create consumer group: %w", err)
		}

		closer.AddNamed("OrderAssembled consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.orderAssembledConsumerGroup = group
	}

	return d.orderAssembledConsumerGroup, nil
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

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) (kafka.Consumer, error) {
	if d.orderAssembledConsumer == nil {
		group, err := d.OrderAssembledConsumerGroup(ctx)
		if err != nil {
			return nil, err
		}

		cfg := config.AppConfig()
		topics := []string{cfg.OrderAssembledConsumer.Topic()}

		d.orderAssembledConsumer = kafkaConsumer.NewConsumer(
			group,
			topics,
			logger.Logger(),
		)
	}

	return d.orderAssembledConsumer, nil
}

func (d *diContainer) OrderPaidDecoder(_ context.Context) (kafkaConverter.OrderPaidDecoder, error) {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder, nil
}

func (d *diContainer) OrderAssembledDecoder(_ context.Context) (kafkaConverter.OrderAssembledDecoder, error) {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder, nil
}

func (d *diContainer) TelegramClient(_ context.Context) (*telegram.Client, error) {
	if d.telegramClient == nil {
		cfg := config.AppConfig()

		client, err := telegram.NewClient(cfg.TelegramBot.Token())
		if err != nil {
			return nil, fmt.Errorf("create telegram client: %w", err)
		}

		client.RegisterStartHandler()

		d.telegramClient = client
	}

	return d.telegramClient, nil
}

func (d *diContainer) TelegramService(ctx context.Context) (svc.TelegramService, error) {
	if d.telegramService == nil {
		client, err := d.TelegramClient(ctx)
		if err != nil {
			return nil, err
		}

		cfg := config.AppConfig()

		service, err := telegramSvc.NewService(client, cfg.TelegramBot.ChatID())
		if err != nil {
			return nil, fmt.Errorf("create telegram service: %w", err)
		}

		d.telegramService = service
	}

	return d.telegramService, nil
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) (svc.OrderPaidConsumerService, error) {
	if d.orderPaidConsumerService == nil {
		consumer, err := d.OrderPaidConsumer(ctx)
		if err != nil {
			return nil, err
		}

		decoder, err := d.OrderPaidDecoder(ctx)
		if err != nil {
			return nil, err
		}

		telegramService, err := d.TelegramService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderPaidConsumerService = orderPaidConsumerSvc.NewService(
			consumer,
			decoder,
			telegramService,
		)
	}

	return d.orderPaidConsumerService, nil
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) (svc.OrderAssembledConsumerService, error) {
	if d.orderAssembledConsumerService == nil {
		consumer, err := d.OrderAssembledConsumer(ctx)
		if err != nil {
			return nil, err
		}

		decoder, err := d.OrderAssembledDecoder(ctx)
		if err != nil {
			return nil, err
		}

		telegramService, err := d.TelegramService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderAssembledConsumerService = orderAssembledConsumerSvc.NewService(
			consumer,
			decoder,
			telegramService,
		)
	}

	return d.orderAssembledConsumerService, nil
}

func (d *diContainer) API(ctx context.Context) (*v1.API, error) {
	if d.api == nil {
		telegramClient, err := d.TelegramClient(ctx)
		if err != nil {
			return nil, err
		}

		d.api = v1.NewAPI(telegramClient)
	}

	return d.api, nil
}
