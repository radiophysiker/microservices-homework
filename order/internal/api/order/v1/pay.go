package v1

import (
	"context"
	"errors"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// PayOrder проводит оплату заказа
func (a *API) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	paymentMethod := converter.ToModelPaymentMethod(req.GetPaymentMethod())

	order, err := a.orderService.PayOrder(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderv1.NotFoundError{
				Error:   orderv1.NotFoundErrorErrorNotFound,
				Message: "order not found",
			}, nil
		case errors.Is(err, model.ErrOrderCannotBePaid):
			return &orderv1.ConflictError{
				Error:   orderv1.ConflictErrorErrorConflict,
				Message: "Order cannot be paid",
			}, nil
		default:
			return &orderv1.InternalServerError{
				Error:   orderv1.InternalServerErrorErrorInternalServerError,
				Message: "failed to pay order",
			}, nil
		}
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: *order.TransactionUUID,
	}, nil
}
