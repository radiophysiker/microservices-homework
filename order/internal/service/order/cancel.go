package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// CancelOrder отменяет заказ
func (s *Service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == model.StatusPaid {
		return nil, model.ErrOrderCannotBeCancelled
	}

	order.Status = model.StatusCancelled

	updated, err := s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return updated, nil
}
