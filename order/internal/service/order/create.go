package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// CreateOrder создает новый заказ
func (s *Service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.NewInvalidOrderDataError("part UUIDs cannot be empty")
	}

	partUUIDStrings := make([]string, len(partUUIDs))
	for i, partUUID := range partUUIDs {
		partUUIDStrings[i] = partUUID.String()
	}

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDStrings)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInventoryServiceUnavailable, err)
	}

	if len(parts) != len(partUUIDs) {
		return nil, model.NewInvalidOrderDataError("some parts not found")
	}

	var totalPrice float64
	items := make([]model.OrderItem, 0, len(parts))

	for _, part := range parts {
		totalPrice += part.Price

		parsedPartUUID, parseErr := uuid.Parse(part.UUID)
		if parseErr != nil {
			return nil, model.NewInvalidOrderDataError("invalid part UUID: " + part.UUID)
		}

		items = append(items, model.OrderItem{
			PartUUID: parsedPartUUID,
			Quantity: 1,
			Price:    &part.Price,
		})
	}

	order := &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
