package grpc

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// InventoryClient представляет интерфейс для работы с inventory service
type InventoryClient interface {
	// ListParts возвращает список деталей по UUID
	ListParts(ctx context.Context, partUUIDs []string) ([]*model.Part, error)
}

// PaymentClient представляет интерфейс для работы с payment service
type PaymentClient interface {
	// PayOrder проводит оплату заказа
	PayOrder(ctx context.Context, userUUID, orderUUID string, paymentMethod paymentpb.PaymentMethod) (string, error)
}
