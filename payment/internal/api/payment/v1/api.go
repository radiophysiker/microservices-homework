package v1

import (
	"github.com/radiophysiker/microservices-homework/payment/internal/service"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// API представляет API слой для payment service
type API struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

// NewAPI создает новый экземпляр API
func NewAPI(paymentService service.PaymentService) *API {
	return &API{
		paymentService: paymentService,
	}
}
