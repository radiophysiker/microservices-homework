package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/radiophysiker/microservices-homework/order/internal/config"
)

// Connect устанавливает пул соединений с PostgresSQL на основе переданной конфигурации cfg.
// Функция настраивает параметры пула (MaxConns, MinConns, MaxConnLifetime, MaxConnIdleTime),
// выполняет подключение и проверку (Ping). При неудаче возвращается ошибка.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	pc, err := pgxpool.ParseConfig(config.AppConfig().Postgres.DSN())
	if err != nil {
		return nil, err
	}

	pc.MaxConns = config.AppConfig().Postgres.PoolMaxConns()
	pc.MinConns = config.AppConfig().Postgres.PoolMinConns()
	pc.MaxConnLifetime = config.AppConfig().Postgres.PoolMaxConnLifetime()
	pc.MaxConnIdleTime = config.AppConfig().Postgres.PoolMaxConnIdleTime()

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
