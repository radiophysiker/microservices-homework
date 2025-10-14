package order

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// GetOrder возвращает заказ по UUID
func (r *Repository) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, exists := r.orders[orderUUID]
	if !exists {
		return nil, model.NewOrderNotFoundError(orderUUID)
	}

	return converter.ToServiceOrder(repoOrder), nil
}
