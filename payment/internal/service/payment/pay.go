package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/payment/internal/model"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа
func (s *Service) PayOrder(ctx context.Context, userUUID, orderUUID string, paymentMethod pb.PaymentMethod) (string, error) {
	if userUUID == "" {
		return "", fmt.Errorf("%w: user_uuid is empty", model.ErrInvalidPaymentRequest)
	}

	if orderUUID == "" {
		return "", fmt.Errorf("%w: order_uuid is empty", model.ErrInvalidPaymentRequest)
	}

	if paymentMethod == pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return "", fmt.Errorf("%w: unspecified payment method", model.ErrInvalidPaymentRequest)
	}

	transactionUUID := uuid.New().String()

	logger.Info(ctx, "Оплата прошла успешно",
		zap.String("user_uuid", userUUID),
		zap.String("order_uuid", orderUUID),
		zap.String("payment_method", paymentMethod.String()),
		zap.String("transaction_uuid", transactionUUID))

	return transactionUUID, nil
}
