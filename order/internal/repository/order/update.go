package order

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// UpdateOrder обновляет заказ
func (r *Repository) UpdateOrder(_ context.Context, order *model.Order) error {
	repoOrder := converter.ToRepoOrder(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.orders[order.OrderUUID.String()]
	if !exists {
		return model.NewOrderNotFoundError(order.OrderUUID.String())
	}

	r.orders[order.OrderUUID.String()] = repoOrder

	return nil
}
