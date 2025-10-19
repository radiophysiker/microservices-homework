package order

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// CreateOrder создает новый заказ
func (r *Repository) CreateOrder(_ context.Context, order *model.Order) error {
	repoOrder := converter.ToRepoOrder(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID.String()] = repoOrder

	return nil
}
