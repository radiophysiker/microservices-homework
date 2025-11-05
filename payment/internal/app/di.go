package app

import (
	apiv1 "github.com/radiophysiker/microservices-homework/payment/internal/api/payment/v1"
	"github.com/radiophysiker/microservices-homework/payment/internal/service"
	paymentSvc "github.com/radiophysiker/microservices-homework/payment/internal/service/payment"
)

type diContainer struct {
	paymentService service.PaymentService
	api            *apiv1.API
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentService() service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentSvc.NewService()
	}

	return d.paymentService
}

func (d *diContainer) API() *apiv1.API {
	if d.api == nil {
		d.api = apiv1.NewAPI(d.PaymentService())
	}

	return d.api
}
