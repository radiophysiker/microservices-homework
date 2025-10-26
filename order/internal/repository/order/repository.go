package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository реализует интерфейс OrderRepository
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository создает новый экземпляр Repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
