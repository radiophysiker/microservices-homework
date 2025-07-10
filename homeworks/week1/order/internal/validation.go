package internal

import (
	"fmt"

	"github.com/google/uuid"
	orderv1 "github.com/radiophysiker/microservices-homework/week1/shared/pkg/openapi/order/v1"
	inventorypb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/inventory/v1"
)

// Order представляет заказ в системе
type Order struct {
	OrderUUID       uuid.UUID                      `json:"order_uuid"`
	UserUUID        uuid.UUID                      `json:"user_uuid"`
	PartUUIDs       []uuid.UUID                    `json:"part_uuids"`
	TotalPrice      float64                        `json:"total_price"`
	TransactionUUID *uuid.UUID                     `json:"transaction_uuid,omitempty"`
	PaymentMethod   *orderv1.OrderDtoPaymentMethod `json:"payment_method,omitempty"`
	Status          orderv1.OrderStatus            `json:"status"`
}

// ValidateOrderExists проверяет существование заказа
func ValidateOrderExists(order *Order) error {
	if order == nil {
		return fmt.Errorf("order not found")
	}
	return nil
}

// ValidateOrderCanBePaid проверяет, может ли заказ быть оплачен
func ValidateOrderCanBePaid(order *Order) error {
	if order.Status != orderv1.OrderStatusPENDINGPAYMENT {
		return fmt.Errorf("order cannot be paid")
	}
	return nil
}

// ValidateOrderCanBeCancelled проверяет, может ли заказ быть отменен
func ValidateOrderCanBeCancelled(order *Order) error {
	if order.Status == orderv1.OrderStatusPAID {
		return fmt.Errorf("paid order cannot be cancelled")
	}
	return nil
}

// ValidateAllPartsFound проверяет, что все запрошенные детали найдены
func ValidateAllPartsFound(requestedPartsCount int, foundParts []*inventorypb.Part) error {
	if len(foundParts) != requestedPartsCount {
		return fmt.Errorf("some parts not found")
	}
	return nil
}

// ValidateTransactionUUID проверяет корректность UUID транзакции
func ValidateTransactionUUID(transactionUUIDStr string) (uuid.UUID, error) {
	transactionUUID, err := uuid.Parse(transactionUUIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid transaction UUID: %w", err)
	}
	return transactionUUID, nil
}

// ValidateCreateOrderRequest проверяет запрос на создание заказа
func ValidateCreateOrderRequest(req *orderv1.CreateOrderRequest) error {
	if len(req.PartUuids) == 0 {
		return fmt.Errorf("part UUIDs cannot be empty")
	}

	// Проверяем что все UUID валидны
	for i, partUUID := range req.PartUuids {
		if partUUID == uuid.Nil {
			return fmt.Errorf("invalid part UUID at index %d", i)
		}
	}

	if req.UserUUID == uuid.Nil {
		return fmt.Errorf("user UUID cannot be empty")
	}

	return nil
}

// ValidatePayOrderRequest проверяет запрос на оплату заказа
func ValidatePayOrderRequest(req *orderv1.PayOrderRequest) error {
	if req.PaymentMethod == orderv1.PaymentMethodUNKNOWN {
		return fmt.Errorf("payment method cannot be UNKNOWN")
	}
	return nil
}

// CalculateTotalPrice вычисляет общую стоимость заказа
func CalculateTotalPrice(parts []*inventorypb.Part) float64 {
	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}
	return totalPrice
}
