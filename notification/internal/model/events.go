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

// OrderPaid представляет событие об оплате заказа
type OrderPaid struct {
	EventUUID       uuid.UUID
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
}

// ShipAssembled представляет событие о завершении сборки корабля
type ShipAssembled struct {
	EventUUID    uuid.UUID
	OrderUUID    uuid.UUID
	UserUUID     uuid.UUID
	BuildTimeSec int64
}
