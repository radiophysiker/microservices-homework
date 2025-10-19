package v1

import (
	"context"
	"errors"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// GetOrder возвращает заказ по UUID
func (a *API) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Error:   orderv1.NotFoundErrorErrorNotFound,
				Message: "order not found",
			}, nil
		}

		return &orderv1.InternalServerError{
			Error:   orderv1.InternalServerErrorErrorInternalServerError,
			Message: "failed to get order",
		}, nil
	}

	return converter.ToOrderDto(order), nil
}
