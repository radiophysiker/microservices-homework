package session

import (
	"time"

	"github.com/radiophysiker/microservices-homework/platform/pkg/cache"
)

// Repository реализует интерфейс SessionRepository для работы с сессиями в Redis.
type Repository struct {
	client cache.RedisClient
	ttl    time.Duration
}

// NewRepository создает новый экземпляр Repository для работы с сессиями.
// Принимает Redis клиент и TTL для сессий.
func NewRepository(client cache.RedisClient, ttl time.Duration) *Repository {
	return &Repository{
		client: client,
		ttl:    ttl,
	}
}
