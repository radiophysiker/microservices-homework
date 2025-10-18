package order

import (
	"sync"

	repoModel "github.com/radiophysiker/microservices-homework/order/internal/repository/model"
)

// Repository реализует интерфейс OrderRepository
type Repository struct {
	mu     sync.RWMutex
	orders map[string]*repoModel.Order
}

// NewRepository создает новый экземпляр Repository
func NewRepository() *Repository {
	return &Repository{
		orders: make(map[string]*repoModel.Order),
	}
}

// GetOrderCount возвращает количество заказов в репозитории (для тестирования)
func (r *Repository) GetOrderCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.orders)
}
