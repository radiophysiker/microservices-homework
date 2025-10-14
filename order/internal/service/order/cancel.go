package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// CancelOrder отменяет заказ
func (s *Service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == orderv1.OrderStatusPAID {
		return nil, model.ErrOrderCannotBeCancelled
	}

	order.Status = orderv1.OrderStatusCANCELLED

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}
