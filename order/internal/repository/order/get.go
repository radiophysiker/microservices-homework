package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// GetOrder возвращает заказ по UUID
func (r *Repository) GetOrder(_ context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, err := uuid.Parse(orderUUID)
	if err != nil {
		return nil, model.NewInvalidOrderDataError(orderUUID)
	}

	repoOrder, exists := r.orders[orderUUID]
	if !exists {
		return nil, model.NewOrderNotFoundError(orderUUID)
	}

	return converter.ToServiceOrder(repoOrder), nil
}
