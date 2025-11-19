package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              string        `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT" envDefault:"5s"`
	MaxIdle           int           `env:"REDIS_MAX_IDLE" envDefault:"10"`
	IdleTimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT" envDefault:"5m"`
}

type RedisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*RedisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &RedisConfig{raw: raw}, nil
}

func (cfg *RedisConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *RedisConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *RedisConfig) ConnectionTimeout() time.Duration {
	return cfg.raw.ConnectionTimeout
}

func (cfg *RedisConfig) MaxIdle() int {
	return cfg.raw.MaxIdle
}

func (cfg *RedisConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}
