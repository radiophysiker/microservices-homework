package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// DB
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string
	DBAuthDB string

	// Ports
	GRPCAddr string
}

// Load загружает конфигурацию из переменных окружения и возвращает заполненную Config.
// Возвращает ошибку, если какие-либо значения окружения имеют некорректный формат.
func Load() (Config, error) {
	port, err := strconv.Atoi(getEnv("INVENTORY_DB_PORT", "27017"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid INVENTORY_DB_PORT: %w", err)
	}

	return Config{
		DBHost:   getEnv("INVENTORY_DB_HOST", "localhost"),
		DBPort:   port,
		DBUser:   getEnv("INVENTORY_DB_USER", "inventory-service-user"),
		DBPass:   getEnv("INVENTORY_DB_PASSWORD", ""),
		DBName:   getEnv("INVENTORY_DB_NAME", "inventory-service"),
		DBAuthDB: getEnv("INVENTORY_DB_AUTH_DB", "admin"),
		GRPCAddr: getEnv("INVENTORY_GRPC_ADDR", "localhost:50051"),
	}, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
