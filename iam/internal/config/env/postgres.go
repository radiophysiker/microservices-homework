package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Database string `env:"POSTGRES_DB,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`

	MaxConns           int32         `env:"POSTGRES_MAX_CONNS" envDefault:"10"`
	MinConns           int32         `env:"POSTGRES_MIN_CONNS" envDefault:"2"`
	MaxConnLifetime    time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime    time.Duration `env:"POSTGRES_MAX_CONN_IDLE" envDefault:"30m"`
	HealthCheckPeriod  time.Duration `env:"POSTGRES_HEALTH_CHECK_PERIOD" envDefault:"1m"`
	MaxConnLifetimeJit time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME_JITTER" envDefault:"0s"`
}

type PostgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*PostgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &PostgresConfig{raw: raw}, nil
}

func (cfg *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.SSLMode,
	)
}

func (cfg *PostgresConfig) PoolMaxConns() int32 {
	return cfg.raw.MaxConns
}

func (cfg *PostgresConfig) PoolMinConns() int32 {
	return cfg.raw.MinConns
}

func (cfg *PostgresConfig) PoolMaxConnLifetime() time.Duration {
	return cfg.raw.MaxConnLifetime
}

func (cfg *PostgresConfig) PoolMaxConnIdleTime() time.Duration {
	return cfg.raw.MaxConnIdleTime
}
