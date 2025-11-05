package migrator

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

// Run выполняет миграции базы данных используя goose.
// pool - пул подключений к PostgreSQL
// dir - директория с SQL файлами миграций
func Run(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error(ctx, "failed to close SQL DB", zap.Error(err))
		}
	}()

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	return goose.UpContext(ctx, sqlDB, dir)
}
