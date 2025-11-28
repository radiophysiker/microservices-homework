package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

// PayOrder проводит оплату заказа
func (s *Service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (*model.Order, error) {
	// Получаем заказ
	order, err := s.orderRepository.GetOrder(ctx, orderUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == model.StatusPaid {
		return nil, model.ErrOrderCannotBePaid
	}

	if order.Status == model.StatusCancelled {
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

	updated, err := s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	eventUUID := uuid.New()
	orderPaidEvent := model.OrderPaid{
		EventUUID:       eventUUID,
		OrderUUID:       updated.OrderUUID,
		UserUUID:        updated.UserUUID,
		PaymentMethod:   *updated.PaymentMethod,
		TransactionUUID: *updated.TransactionUUID,
	}

	if err := s.orderProducer.ProduceOrderPaid(ctx, orderPaidEvent); err != nil {
		logger.Error(ctx, "Failed to publish OrderPaid event",
			zap.Error(err),
			zap.String("order_uuid", updated.OrderUUID.String()),
			zap.String("event_uuid", eventUUID.String()),
		)
	}

	if s.revenueCounter != nil {
		s.revenueCounter.Add(ctx, updated.TotalPrice)
	}

	return updated, nil
}
