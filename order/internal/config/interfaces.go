package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	InventoryAddress() string
}

type PostgresConfig interface {
	DSN() string
	PoolMaxConns() int32
	PoolMinConns() int32
	PoolMaxConnLifetime() time.Duration
	PoolMaxConnIdleTime() time.Duration
}

type PaymentGRPCConfig interface {
	PaymentAddress() string
}

type IAMGRPCConfig interface {
	IAMAddress() string
}

type OrderGRPCConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
}

type MigrationsConfig interface {
	Directory() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
