package model

import (
	"errors"
	"fmt"
)

var (
	// ErrOrderNotFound - ошибка "заказ не найден"
	ErrOrderNotFound = errors.New("order not found")
	// ErrInvalidOrderData - ошибка "некорректные данные заказа"
	ErrInvalidOrderData = errors.New("invalid order data")
	// ErrOrderCannotBePaid - ошибка "заказ не может быть оплачен"
	ErrOrderCannotBePaid = errors.New("order cannot be paid")
	// ErrOrderCannotBeCancelled - ошибка "заказ не может быть отменен"
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
	// ErrPartNotFound - ошибка "деталь не найдена"
	ErrPartNotFound = errors.New("part not found")
	// ErrInventoryServiceUnavailable - ошибка "сервис инвентаря недоступен"
	ErrInventoryServiceUnavailable = errors.New("inventory service unavailable")
	// ErrPaymentServiceUnavailable - ошибка "сервис платежей недоступен"
	ErrPaymentServiceUnavailable = errors.New("payment service unavailable")
)

// NewOrderNotFoundError создает ошибку "заказ не найден"
func NewOrderNotFoundError(orderUUID string) error {
	return fmt.Errorf("%w: %s", ErrOrderNotFound, orderUUID)
}

// NewInvalidOrderDataError создает ошибку "некорректные данные заказа"
func NewInvalidOrderDataError(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidOrderData, msg)
}
