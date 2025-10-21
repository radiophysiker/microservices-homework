package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// GetOrder возвращает заказ по UUID
func (s *Service) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}
