package order

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	clientGrpc "github.com/radiophysiker/microservices-homework/order/internal/client/grpc"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
	orderProducer "github.com/radiophysiker/microservices-homework/order/internal/service/producer/order_producer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

// Service реализует интерфейс OrderService
type Service struct {
	orderRepository repository.OrderRepository
	inventoryClient clientGrpc.InventoryClient
	paymentClient   clientGrpc.PaymentClient
	orderProducer   *orderProducer.Service
	ordersCounter   metric.Int64Counter
	revenueCounter  metric.Float64Counter
}

// NewService создает новый экземпляр Service
func NewService(
	ctx context.Context,
	orderRepository repository.OrderRepository,
	inventoryClient clientGrpc.InventoryClient,
	paymentClient clientGrpc.PaymentClient,
	orderProducer *orderProducer.Service,
) *Service {
	meter := otel.Meter("order-service")

	ordersCounter, err := meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of orders created"),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create orders_total counter", zap.Error(err))
	}

	revenueCounter, err := meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from orders"),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create orders_revenue_total counter", zap.Error(err))
	}

	return &Service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
		ordersCounter:   ordersCounter,
		revenueCounter:  revenueCounter,
	}
}
