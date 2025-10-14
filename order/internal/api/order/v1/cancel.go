package v1

import (
	"context"
	"errors"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// CancelOrder отменяет заказ
func (a *API) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	_, err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderv1.NotFoundError{
				Error:   orderv1.NotFoundErrorErrorNotFound,
				Message: "order not found",
			}, nil
		case errors.Is(err, model.ErrOrderCannotBeCancelled):
			return &orderv1.ConflictError{
				Error:   orderv1.ConflictErrorErrorConflict,
				Message: "Order cannot be cancelled",
			}, nil
		default:
			return &orderv1.InternalServerError{
				Error:   orderv1.InternalServerErrorErrorInternalServerError,
				Message: "failed to cancel order",
			}, nil
		}
	}

	return &orderv1.CancelOrderNoContent{}, nil
}
