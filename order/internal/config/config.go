package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/radiophysiker/microservices-homework/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	OrderGRPC     OrderGRPCConfig
	OrderHTTP     OrderHTTPConfig
	Postgres      PostgresConfig
	Migrations    MigrationsConfig
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

	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
		OrderGRPC:     orderGRPCCfg,
		OrderHTTP:     httpCfg,
		Postgres:      postgresCfg,
		Migrations:    migrationsCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
