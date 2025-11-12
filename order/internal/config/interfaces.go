package config

import (
	"time"
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

type OrderGRPCConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
}

type MigrationsConfig interface {
	Directory() string
}
