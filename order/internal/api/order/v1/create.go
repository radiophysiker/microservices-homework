package v1

import (
	"context"
	"errors"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// CreateOrder создает новый заказ
func (a *API) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	if len(req.GetPartUuids()) == 0 {
		return &orderv1.BadRequestError{
			Error:   orderv1.BadRequestErrorErrorBadRequest,
			Message: "part UUIDs cannot be empty",
		}, nil
	}

	order, err := a.orderService.CreateOrder(ctx, req.GetUserUUID(), req.GetPartUuids())
	if err != nil {
		if errors.Is(err, model.ErrInvalidOrderData) {
			return &orderv1.BadRequestError{
				Error:   orderv1.BadRequestErrorErrorBadRequest,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrInventoryServiceUnavailable) {
			return &orderv1.InternalServerError{
				Error:   orderv1.InternalServerErrorErrorInternalServerError,
				Message: "inventory service unavailable",
			}, nil
		}

		return &orderv1.InternalServerError{
			Error:   orderv1.InternalServerErrorErrorInternalServerError,
			Message: "failed to create order",
		}, nil
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
