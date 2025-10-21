package v1

import (
	"github.com/radiophysiker/microservices-homework/order/internal/service"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

// API представляет API слой для order service
type API struct {
	orderpb.UnimplementedOrderServiceServer
	orderService service.OrderService
}

// NewAPI создает новый экземпляр API
func NewAPI(orderService service.OrderService) *API {
	return &API{
		orderService: orderService,
	}
}
