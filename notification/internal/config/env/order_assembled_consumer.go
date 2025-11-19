//nolint:dupl // Файл похож на order_paid_consumer.go, но это разные конфигурации для разных топиков
package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type OrderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	raw OrderAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw OrderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{raw: raw}, nil
}

func (cfg *orderAssembledConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderAssembledConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
