package order

import (
	clientGrpc "github.com/radiophysiker/microservices-homework/order/internal/client/grpc"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
	orderProducer "github.com/radiophysiker/microservices-homework/order/internal/service/producer/order_producer"
)

// Service реализует интерфейс OrderService
type Service struct {
	orderRepository repository.OrderRepository
	inventoryClient clientGrpc.InventoryClient
	paymentClient   clientGrpc.PaymentClient
	orderProducer   *orderProducer.Service
}

// NewService создает новый экземпляр Service
func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient clientGrpc.InventoryClient,
	paymentClient clientGrpc.PaymentClient,
	orderProducer *orderProducer.Service,
) *Service {
	return &Service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
	}
}
