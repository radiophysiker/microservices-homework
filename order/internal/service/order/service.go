package order

import (
	clientGrpc "github.com/radiophysiker/microservices-homework/order/internal/client/grpc"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
)

// Service реализует интерфейс OrderService
type Service struct {
	orderRepository repository.OrderRepository
	inventoryClient clientGrpc.InventoryClient
	paymentClient   clientGrpc.PaymentClient
}

// NewService создает новый экземпляр Service
func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient clientGrpc.InventoryClient,
	paymentClient clientGrpc.PaymentClient,
) *Service {
	return &Service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
