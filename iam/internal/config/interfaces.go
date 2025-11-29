package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type PostgresConfig interface {
	DSN() string
	PoolMaxConns() int32
	PoolMinConns() int32
	PoolMaxConnLifetime() time.Duration
	PoolMaxConnIdleTime() time.Duration
}

type MigrationsConfig interface {
	Directory() string
}

type RedisConfig interface {
	Host() string
	Port() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type IAMGRPCConfig interface {
	Host() string
	Port() string
	Address() string
}

type SessionConfig interface {
	TTL() time.Duration
}
