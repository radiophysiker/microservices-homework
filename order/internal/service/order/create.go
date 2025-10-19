package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// CreateOrder создает новый заказ
func (s *Service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	// Валидация входных данных
	if len(partUUIDs) == 0 {
		return nil, model.NewInvalidOrderDataError("part UUIDs cannot be empty")
	}

	// Конвертируем UUID в строки для вызова inventory service
	partUUIDStrings := make([]string, len(partUUIDs))
	for i, partUUID := range partUUIDs {
		partUUIDStrings[i] = partUUID.String()
	}

	// Получаем информацию о деталях из inventory service
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDStrings)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInventoryServiceUnavailable, err)
	}

	// Проверяем, что все запрошенные детали найдены
	if len(parts) != len(partUUIDs) {
		return nil, model.NewInvalidOrderDataError("some parts not found")
	}

	// Рассчитываем общую стоимость
	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	// Создаем заказ
	order := &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     orderv1.OrderStatusPENDINGPAYMENT,
	}

	// Сохраняем заказ в repository
	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
