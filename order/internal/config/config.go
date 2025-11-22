package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/radiophysiker/microservices-homework/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	InventoryGRPC          InventoryGRPCConfig
	PaymentGRPC            PaymentGRPCConfig
	IAMGRPC                IAMGRPCConfig
	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderGRPC              OrderGRPCConfig
	OrderHTTP              OrderHTTPConfig
	Postgres               PostgresConfig
	Migrations             MigrationsConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	iamGRPCCfg, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	orderGRPCCfg, err := env.NewOrderGRPCConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	migrationsCfg, err := env.NewMigrationsConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		InventoryGRPC:          inventoryGRPCCfg,
		PaymentGRPC:            paymentGRPCCfg,
		IAMGRPC:                iamGRPCCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		OrderGRPC:              orderGRPCCfg,
		OrderHTTP:              httpCfg,
		Postgres:               postgresCfg,
		Migrations:             migrationsCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
