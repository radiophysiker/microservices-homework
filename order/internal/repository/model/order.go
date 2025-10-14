package model

import (
	"github.com/google/uuid"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// Order представляет заказ в repository слое
type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *orderv1.OrderDtoPaymentMethod
	Status          orderv1.OrderStatus
}
