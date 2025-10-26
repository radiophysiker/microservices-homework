package migrator

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Run(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close SQL DB: %v", err)
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
