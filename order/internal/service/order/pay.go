package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// PayOrder проводит оплату заказа
func (s *Service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (*model.Order, error) {
	// Получаем заказ
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Проверяем, что заказ может быть оплачен
	if order.Status != model.StatusPendingPayment {
		return nil, model.ErrOrderCannotBePaid
	}

	// Проводим оплату через payment service
	transactionUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.UserUUID.String(),
		order.OrderUUID.String(),
		converter.PaymentMethodToProtobuf(paymentMethod),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrPaymentServiceUnavailable, err)
	}

	// Обновляем заказ
	parsedTransactionUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction UUID: %w", err)
	}

	order.TransactionUUID = &parsedTransactionUUID
	order.PaymentMethod = &paymentMethod
	order.Status = model.StatusPaid

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}
