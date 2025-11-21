package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// OrderService представляет интерфейс для работы с заказами
type OrderService interface {
	// CreateOrder создает новый заказ
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error)
	// GetOrder возвращает заказ по UUID
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	// PayOrder проводит оплату заказа
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (*model.Order, error)
	// CancelOrder отменяет заказ
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
}

// OrderConsumerService представляет интерфейс для consumer'а событий ShipAssembled
type OrderConsumerService interface {
	// RunConsumer запускает consumer для обработки событий ShipAssembled
	RunConsumer(ctx context.Context) error
}
