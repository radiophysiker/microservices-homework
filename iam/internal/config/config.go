package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/radiophysiker/microservices-homework/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	Postgres   PostgresConfig
	Migrations MigrationsConfig
	Redis      RedisConfig
	IAMGRPC    IAMGRPCConfig
	Session    SessionConfig
}

// Load загружает конфигурацию из переменных окружения.
// Принимает опциональные пути к .env файлам.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
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

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	iamGRPCCfg, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerCfg,
		Postgres:   postgresCfg,
		Migrations: migrationsCfg,
		Redis:      redisCfg,
		IAMGRPC:    iamGRPCCfg,
		Session:    sessionCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию приложения.
// Должен быть вызван после Load.
func AppConfig() *config {
	return appConfig
}
