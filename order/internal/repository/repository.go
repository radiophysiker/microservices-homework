package repository

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// OrderRepository представляет интерфейс для работы с заказами в repository слое
type OrderRepository interface {
	// CreateOrder создает новый заказ
	CreateOrder(ctx context.Context, order *model.Order) error
	// GetOrder возвращает заказ по UUID
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	// UpdateOrder обновляет заказ и возвращает актуальное состояние
	UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
}
