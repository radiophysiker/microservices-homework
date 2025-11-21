package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type OrderPaidProducerEnvConfig struct {
	Topic string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	raw OrderPaidProducerEnvConfig
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw OrderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{raw: raw}, nil
}

func (cfg *orderPaidProducerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
