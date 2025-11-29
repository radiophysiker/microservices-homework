package app

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/radiophysiker/microservices-homework/order/internal/api/order/v1"
	clientGrpc "github.com/radiophysiker/microservices-homework/order/internal/client/grpc"
	inventoryClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/payment/v1"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/order/internal/converter/kafka/decoder"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
	orderRepo "github.com/radiophysiker/microservices-homework/order/internal/repository/order"
	"github.com/radiophysiker/microservices-homework/order/internal/service"
	orderConsumerSvc "github.com/radiophysiker/microservices-homework/order/internal/service/consumer/order_consumer"
	orderSvc "github.com/radiophysiker/microservices-homework/order/internal/service/order"
	orderProducerSvc "github.com/radiophysiker/microservices-homework/order/internal/service/producer/order_producer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/kafka"
	kafkaConsumer "github.com/radiophysiker/microservices-homework/platform/pkg/kafka/consumer"
	kafkaProducer "github.com/radiophysiker/microservices-homework/platform/pkg/kafka/producer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	grpcMiddleware "github.com/radiophysiker/microservices-homework/platform/pkg/middleware/grpc"
	httpMiddleware "github.com/radiophysiker/microservices-homework/platform/pkg/middleware/http"
	"github.com/radiophysiker/microservices-homework/platform/pkg/tracing"
	authpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/auth/v1"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	pool            *pgxpool.Pool
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
	iamConn         *grpc.ClientConn
	orderRepository repository.OrderRepository
	inventoryClient clientGrpc.InventoryClient
	paymentClient   clientGrpc.PaymentClient
	iamClient       authpb.AuthServiceClient
	orderService    service.OrderService
	api             *apiv1.API

	orderPaidSyncProducer sarama.SyncProducer
	orderPaidProducer     kafka.Producer
	orderProducerService  *orderProducerSvc.Service

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      kafka.Consumer
	orderAssembledDecoder       *decoder.Decoder
	orderConsumerService        service.OrderConsumerService
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Pool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pool == nil {
		pc, err := pgxpool.ParseConfig(config.AppConfig().Postgres.DSN())
		if err != nil {
			return nil, fmt.Errorf("parse postgres config: %w", err)
		}

		pc.MaxConns = config.AppConfig().Postgres.PoolMaxConns()
		pc.MinConns = config.AppConfig().Postgres.PoolMinConns()
		pc.MaxConnLifetime = config.AppConfig().Postgres.PoolMaxConnLifetime()
		pc.MaxConnIdleTime = config.AppConfig().Postgres.PoolMaxConnIdleTime()

		ctxConnect, cancelConnect := context.WithTimeout(ctx, 10*time.Second)
		defer cancelConnect()

		pool, err := pgxpool.NewWithConfig(ctxConnect, pc)
		if err != nil {
			return nil, fmt.Errorf("create postgres pool: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, fmt.Errorf("ping postgres: %w", err)
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pool = pool
	}

	return d.pool, nil
}

func (d *diContainer) InventoryConn(ctx context.Context) (*grpc.ClientConn, error) {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.InventoryAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(
				tracing.UnaryClientInterceptor(config.AppConfig().Tracing.ServiceName()),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("connect inventory grpc: %w", err)
		}

		closer.AddNamed("inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn, nil
}

func (d *diContainer) PaymentConn(ctx context.Context) (*grpc.ClientConn, error) {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.PaymentAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(
				tracing.UnaryClientInterceptor(config.AppConfig().Tracing.ServiceName()),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("connect payment grpc: %w", err)
		}

		closer.AddNamed("payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn, nil
}

func (d *diContainer) IAMConn(ctx context.Context) (*grpc.ClientConn, error) {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMGRPC.IAMAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(
				tracing.UnaryClientInterceptor(config.AppConfig().Tracing.ServiceName()),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("connect iam grpc: %w", err)
		}

		closer.AddNamed("iam gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn, nil
}

func (d *diContainer) IAMClient(ctx context.Context) (authpb.AuthServiceClient, error) {
	if d.iamClient == nil {
		conn, err := d.IAMConn(ctx)
		if err != nil {
			return nil, err
		}

		d.iamClient = authpb.NewAuthServiceClient(conn)
	}

	return d.iamClient, nil
}

func (d *diContainer) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if d.orderRepository == nil {
		pool, err := d.Pool(ctx)
		if err != nil {
			return nil, err
		}

		d.orderRepository = orderRepo.NewRepository(pool)
	}

	return d.orderRepository, nil
}

func (d *diContainer) InventoryClient(ctx context.Context) (clientGrpc.InventoryClient, error) {
	if d.inventoryClient == nil {
		conn, err := d.InventoryConn(ctx)
		if err != nil {
			return nil, err
		}

		d.inventoryClient = inventoryClient.NewClient(
			inventorypb.NewInventoryServiceClient(conn),
		)
	}

	return d.inventoryClient, nil
}

func (d *diContainer) PaymentClient(ctx context.Context) (clientGrpc.PaymentClient, error) {
	if d.paymentClient == nil {
		conn, err := d.PaymentConn(ctx)
		if err != nil {
			return nil, err
		}

		d.paymentClient = paymentClient.NewClient(
			paymentpb.NewPaymentServiceClient(conn),
		)
	}

	return d.paymentClient, nil
}

func (d *diContainer) OrderPaidSyncProducer(ctx context.Context) (sarama.SyncProducer, error) {
	if d.orderPaidSyncProducer == nil {
		cfg := config.AppConfig()
		producerCfg := cfg.OrderPaidProducer

		producer, err := sarama.NewSyncProducer(
			cfg.Kafka.Brokers(),
			producerCfg.Config(),
		)
		if err != nil {
			return nil, fmt.Errorf("create sync producer: %w", err)
		}

		closer.AddNamed("OrderPaid sync producer", func(ctx context.Context) error {
			return producer.Close()
		})

		d.orderPaidSyncProducer = producer
	}

	return d.orderPaidSyncProducer, nil
}

func (d *diContainer) OrderPaidProducer(ctx context.Context) (kafka.Producer, error) {
	if d.orderPaidProducer == nil {
		syncProducer, err := d.OrderPaidSyncProducer(ctx)
		if err != nil {
			return nil, err
		}

		cfg := config.AppConfig()
		topic := cfg.OrderPaidProducer.Topic()

		d.orderPaidProducer = kafkaProducer.NewProducer(
			syncProducer,
			topic,
			logger.Logger(),
		)
	}

	return d.orderPaidProducer, nil
}

func (d *diContainer) OrderProducerService(ctx context.Context) (*orderProducerSvc.Service, error) {
	if d.orderProducerService == nil {
		producer, err := d.OrderPaidProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.orderProducerService = orderProducerSvc.NewService(producer)
	}

	return d.orderProducerService, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderService == nil {
		orderRepo, err := d.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}

		inventoryClient, err := d.InventoryClient(ctx)
		if err != nil {
			return nil, err
		}

		paymentClient, err := d.PaymentClient(ctx)
		if err != nil {
			return nil, err
		}

		orderProducer, err := d.OrderProducerService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderService = orderSvc.NewService(
			ctx,
			orderRepo,
			inventoryClient,
			paymentClient,
			orderProducer,
		)
	}

	return d.orderService, nil
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

func (d *diContainer) OrderAssembledDecoder(_ context.Context) (*decoder.Decoder, error) {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder, nil
}

func (d *diContainer) OrderConsumerService(ctx context.Context) (service.OrderConsumerService, error) {
	if d.orderConsumerService == nil {
		consumer, err := d.OrderAssembledConsumer(ctx)
		if err != nil {
			return nil, err
		}

		decoder, err := d.OrderAssembledDecoder(ctx)
		if err != nil {
			return nil, err
		}

		orderRepo, err := d.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.orderConsumerService = orderConsumerSvc.NewService(
			consumer,
			decoder,
			orderRepo,
		)
	}

	return d.orderConsumerService, nil
}

func (d *diContainer) API(ctx context.Context) (*apiv1.API, error) {
	if d.api == nil {
		orderService, err := d.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		d.api = apiv1.NewAPI(orderService)
	}

	return d.api, nil
}

func (d *diContainer) AuthMiddleware(ctx context.Context) (*httpMiddleware.AuthMiddleware, error) {
	iamClient, err := d.IAMClient(ctx)
	if err != nil {
		return nil, err
	}

	return httpMiddleware.NewAuthMiddleware(iamClient), nil
}

func (d *diContainer) AuthInterceptor(ctx context.Context) (*grpcMiddleware.AuthInterceptor, error) {
	iamClient, err := d.IAMClient(ctx)
	if err != nil {
		return nil, err
	}

	return grpcMiddleware.NewAuthInterceptor(iamClient), nil
}
