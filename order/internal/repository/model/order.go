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

// String возвращает строковое представление PaymentMethod
func (pm PaymentMethod) String() string {
	switch pm {
	case PaymentMethodCard:
		return "CARD"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case PaymentMethodInvestorMoney:
		return "INVESTOR_MONEY"
	default:
		return "UNSPECIFIED"
	}
}

// Status представляет статус заказа
type Status int

const (
	StatusUnspecified Status = iota
	StatusPendingPayment
	StatusPaid
	StatusCancelled
)

// String возвращает строковое представление Status
func (s Status) String() string {
	switch s {
	case StatusPendingPayment:
		return "PENDING_PAYMENT"
	case StatusPaid:
		return "PAID"
	case StatusCancelled:
		return "CANCELLED"
	default:
		return "UNSPECIFIED"
	}
}

// Order представляет заказ в repository слое
type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Items           []OrderItem
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          Status
}

// OrderItem представляет позицию заказа в repository слое
type OrderItem struct {
	PartUUID uuid.UUID
	Quantity int
}
