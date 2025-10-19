package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/payment/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа
func (s *Service) PayOrder(_ context.Context, userUUID, orderUUID string, paymentMethod pb.PaymentMethod) (string, error) {
	if userUUID == "" || orderUUID == "" {
		return "", model.ErrInvalidPaymentRequest
	}

	if paymentMethod == pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return "", fmt.Errorf("%w: unspecified payment method", model.ErrInvalidPaymentRequest)
	}

	transactionUUID := uuid.New().String()

	log.Printf("Оплата прошла успешно, user_uuid: %s, order_uuid: %s, payment_method: %s, transaction_uuid: %s",
		userUUID, orderUUID, paymentMethod.String(), transactionUUID)

	return transactionUUID, nil
}
