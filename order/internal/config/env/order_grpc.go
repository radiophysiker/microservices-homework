package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderGRPCEnvConfig struct {
	Host string `env:"ORDER_GRPC_HOST" envDefault:"0.0.0.0"`
	Port string `env:"ORDER_GRPC_PORT" envDefault:"50053"`
}

type orderGRPCConfig struct {
	raw orderGRPCEnvConfig
}

func NewOrderGRPCConfig() (*orderGRPCConfig, error) {
	var raw orderGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderGRPCConfig{raw: raw}, nil
}

func (cfg *orderGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
