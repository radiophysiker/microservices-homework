package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/radiophysiker/microservices-homework/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	TelegramBot            TelegramBotConfig
	HTTP                   HTTPConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	telegramBotCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		TelegramBot:            telegramBotCfg,
		HTTP:                   httpCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
