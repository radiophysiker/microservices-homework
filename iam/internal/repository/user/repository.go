package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository реализует интерфейс UserRepository для работы с пользователями в PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository создает новый экземпляр Repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
