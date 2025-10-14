package service

import (
	"context"

	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// PaymentService представляет интерфейс для работы с платежами
type PaymentService interface {
	// PayOrder проводит оплату заказа
	PayOrder(ctx context.Context, userUUID, orderUUID string, paymentMethod pb.PaymentMethod) (string, error)
}
