package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// OrderService представляет интерфейс для работы с заказами
type OrderService interface {
	// CreateOrder создает новый заказ
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error)
	// GetOrder возвращает заказ по UUID
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	// PayOrder проводит оплату заказа
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod orderv1.PaymentMethod) (*model.Order, error)
	// CancelOrder отменяет заказ
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
}
