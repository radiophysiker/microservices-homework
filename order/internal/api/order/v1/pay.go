package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

// PayOrder проводит оплату заказа
func (a *API) PayOrder(ctx context.Context, req *orderpb.PayOrderRequest) (*orderpb.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order UUID: %v", err)
	}

	paymentMethod := converter.PaymentMethodFromProtobuf(req.PaymentMethod)

	order, err := a.orderService.PayOrder(ctx, orderUUID, paymentMethod)
	if err != nil {
		// Обработка различных типов ошибок
		switch {
		case errors.Is(err, model.ErrInvalidOrderData):
			return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
		case errors.Is(err, model.ErrOrderNotFound):
			return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
		case errors.Is(err, model.ErrOrderCannotBePaid):
			return nil, status.Errorf(codes.FailedPrecondition, "order cannot be paid: %v", err)
		case errors.Is(err, model.ErrPaymentServiceUnavailable):
			return nil, status.Errorf(codes.Unavailable, "payment service unavailable: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to pay order: %v", err)
		}
	}

	if order.TransactionUUID == nil {
		return nil, status.Errorf(codes.Internal, "transaction UUID not set after payment")
	}

	return &orderpb.PayOrderResponse{
		TransactionUuid: order.TransactionUUID.String(),
	}, nil
}
