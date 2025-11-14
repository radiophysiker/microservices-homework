package model

import (
	"github.com/google/uuid"
)

// PaymentMethod представляет способ оплаты
type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

// Status представляет статус заказа
type Status int

const (
	StatusUnspecified Status = iota
	StatusPendingPayment
	StatusPaid
	StatusAssembled
	StatusCancelled
)

// Order представляет заказ в сервисном слое
type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Items           []OrderItem
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          Status
}

type OrderItem struct {
	PartUUID uuid.UUID
	Quantity int
}
