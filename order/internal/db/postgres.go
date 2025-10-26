package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/radiophysiker/microservices-homework/order/internal/config"
)

// Connect устанавливает пул соединений с PostgresSQL на основе переданной конфигурации cfg.
// Функция настраивает параметры пула (MaxConns, MinConns, MaxConnLifetime, MaxConnIdleTime),
// выполняет подключение и проверку (Ping). При неудаче возвращается ошибка.
func Connect(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pc, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pc.MaxConns = cfg.MaxConns
	pc.MinConns = cfg.MinConns
	pc.MaxConnLifetime = cfg.MaxConnLifetime
	pc.MaxConnIdleTime = cfg.MaxConnIdle

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
