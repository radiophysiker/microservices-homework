package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// DB
	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string

	// Ports
	GRPCAddr string // ":50053"
	HTTPAddr string // ":8080"

	// Deps
	InventoryAddr string // "inventory:50051"
	PaymentAddr   string // "payment:50052"

	// Pool
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdle     time.Duration

	// Migrations
	MigrationsDir string
}

// Load загружает конфигурацию из переменных окружения и возвращает заполненную Config.
// Возвращает ошибку, если какие-либо значения окружения имеют некорректный формат
// или отсутствуют обязательные поля.
func Load() (Config, error) {
	var c Config

	var err error

	c.DBHost = getEnv("ORDER_DB_HOST", "localhost")

	c.DBPort, err = parseInt("ORDER_DB_PORT", 5432)
	if err != nil {
		return Config{}, err
	}

	c.DBUser = getEnv("ORDER_DB_USER", "order-service-user")
	c.DBPass = getEnv("ORDER_DB_PASSWORD", "")
	c.DBName = getEnv("ORDER_DB_NAME", "order-service")

	c.GRPCAddr = getEnv("ORDER_GRPC_ADDR", ":50053")
	c.HTTPAddr = getEnv("ORDER_HTTP_ADDR", ":8080")

	c.InventoryAddr = getEnv("INVENTORY_GRPC_ADDR", "localhost:50051")
	c.PaymentAddr = getEnv("PAYMENT_GRPC_ADDR", "localhost:50052")

	c.MaxConns, err = parseInt32("DB_MAX_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	c.MinConns, err = parseInt32("DB_MIN_CONNS", 2)
	if err != nil {
		return Config{}, err
	}

	c.MaxConnLifetime = parseDuration("DB_MAX_CONN_LIFETIME", time.Hour)
	c.MaxConnIdle = parseDuration("DB_MAX_CONN_IDLE", 30*time.Minute)

	c.MigrationsDir = getEnv("MIGRATIONS_DIR", "./migrations")

	if c.DBUser == "" || c.DBPass == "" || c.DBName == "" {
		return Config{}, errors.New("missing required database credentials")
	}

	return c, nil
}

// parseInt32 читает переменную окружения key и пытается распарсить её как int32.
// Если переменная пуста, возвращается значение по умолчанию defaultValue.
// Если значение задано, но не является валидным int32, возвращается ошибка.
func parseInt32(key string, defaultValue int32) (int32, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}

	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid int32 for %s: %w", key, err)
	}

	return int32(n), nil
}

// parseInt читает переменную окружения key и пытается распарсить её как int.
// Если переменная пуста, возвращается значение по умолчанию defaultValue.
// Если значение задано, но не является валидным int, возвращается ошибка.
func parseInt(key string, defaultValue int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid int for %s: %w", key, err)
	}

	return n, nil
}

// getEnv возвращает значение переменной окружения key или defaultValue, если переменная не задана.
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}

// parseDuration читает переменную окружения key и пытается распарсить её как time.Duration.
// Если переменная пуста или парсинг не удался, возвращается значение по умолчанию defaultValue.
func parseDuration(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}

	return d
}
