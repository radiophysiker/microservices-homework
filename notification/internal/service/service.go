package service

import (
	"context"

	"github.com/radiophysiker/microservices-homework/notification/internal/model"
)

// OrderPaidConsumerService представляет интерфейс для consumer'а событий OrderPaid
type OrderPaidConsumerService interface {
	// RunConsumer запускает consumer для обработки событий OrderPaid
	RunConsumer(ctx context.Context) error
}

// OrderAssembledConsumerService представляет интерфейс для consumer'а событий ShipAssembled
type OrderAssembledConsumerService interface {
	// RunConsumer запускает consumer для обработки событий ShipAssembled
	RunConsumer(ctx context.Context) error
}

// TelegramService представляет интерфейс для отправки уведомлений в Telegram
type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event *model.OrderPaid) error
	SendShipAssembledNotification(ctx context.Context, event *model.ShipAssembled) error
	HandleStartCommand(ctx context.Context, chatID string) error
}
