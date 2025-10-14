package v1

import (
	"context"

	"github.com/radiophysiker/microservices-homework/order/internal/service"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// API представляет API слой для order service
type API struct {
	orderService service.OrderService
}

// NewAPI создает новый экземпляр API
func NewAPI(orderService service.OrderService) *API {
	return &API{
		orderService: orderService,
	}
}

// NewError создает стандартизированный ответ с ошибкой для неизвестных ошибок
func (a *API) NewError(ctx context.Context, err error) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderv1.GenericError{
			Error:   "error",
			Message: err.Error(),
		},
	}
}
