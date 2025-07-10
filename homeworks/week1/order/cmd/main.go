package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/radiophysiker/microservices-homework/week1/order/internal"
	orderv1 "github.com/radiophysiker/microservices-homework/week1/shared/pkg/openapi/order/v1"
	inventorypb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/inventory/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server реализует HTTP API для OrderService.
type Server struct {
	orders          map[string]*internal.Order
	mu              sync.RWMutex
	inventoryClient inventorypb.InventoryServiceClient
	paymentClient   paymentpb.PaymentServiceClient
}

// NewServer создает новый экземпляр сервера OrderService.
func NewServer() (*Server, error) {
	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	log.Println("Successfully connected to external services")

	return &Server{
		orders:          make(map[string]*internal.Order),
		inventoryClient: inventorypb.NewInventoryServiceClient(inventoryConn),
		paymentClient:   paymentpb.NewPaymentServiceClient(paymentConn),
	}, nil
}

/*
CreateOrder создает новый заказ.
Валидирует запрос, проверяет наличие деталей в InventoryService,
рассчитывает общую стоимость и сохраняет заказ со статусом PENDING_PAYMENT.
*/
func (s *Server) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	log.Printf("Creating order for user %s with %d parts", req.UserUUID.String(), len(req.PartUuids))
	if err := internal.ValidateCreateOrderRequest(req); err != nil {
		log.Printf("Validation failed: %v", err)
		return internal.NewOrderBadRequestError(err.Error()), nil
	}

	partUUIDs := internal.ConvertUUIDsToStrings(req.PartUuids)
	inventoryResp, err := s.inventoryClient.ListParts(ctx, &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		log.Printf("Failed to get parts from inventory: %v", err)
		return internal.NewOrderInternalError(
			fmt.Sprintf("Failed to get parts from inventory: %v", err),
		), nil
	}

	if err := internal.ValidateAllPartsFound(len(req.PartUuids), inventoryResp.Parts); err != nil {
		log.Printf("Parts validation failed: %v", err)
		return internal.NewOrderBadRequestError(err.Error()), nil
	}

	totalPrice := internal.CalculateTotalPrice(inventoryResp.Parts)

	orderUUID := uuid.New()
	order := &internal.Order{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderv1.OrderStatusPENDINGPAYMENT,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[orderUUID.String()] = order

	log.Printf("Order %s created successfully with total price %.2f", orderUUID.String(), totalPrice)

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// GetOrder возвращает информацию о заказе по его UUID.
func (s *Server) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	orderUUIDStr := params.OrderUUID.String()
	log.Printf("Getting order %s", orderUUIDStr)

	s.mu.RLock()
	defer s.mu.RUnlock()
	order, exists := s.orders[orderUUIDStr]

	if !exists {
		log.Printf("Order %s not found", orderUUIDStr)
		return internal.NewOrderNotFoundError("Order not found"), nil
	}

	log.Printf("Order %s found with status %s", orderUUIDStr, order.Status)

	// Конвертируем внутреннюю структуру в DTO
	return internal.ConvertOrderToDto(order), nil
}

/*
PayOrder проводит оплату заказа.
Валидирует запрос, проверяет возможность оплаты,
вызывает PaymentService и обновляет статус заказа.
*/
func (s *Server) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	orderUUIDStr := params.OrderUUID.String()
	log.Printf("Processing payment for order %s with method %s", orderUUIDStr, req.PaymentMethod)

	if err := internal.ValidatePayOrderRequest(req); err != nil {
		log.Printf("Payment request validation failed: %v", err)
		return internal.NewPayOrderInternalError(err.Error()), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderUUIDStr]
	if !exists {
		log.Printf("Order %s not found for payment", orderUUIDStr)
		return internal.NewPayOrderNotFoundError("Order not found"), nil
	}

	if err := internal.ValidateOrderCanBePaid(order); err != nil {
		log.Printf("Order %s cannot be paid: %v", orderUUIDStr, err)
		return internal.NewPayOrderConflictError(err.Error()), nil
	}

	paymentResp, err := s.paymentClient.PayOrder(ctx, &paymentpb.PayOrderRequest{
		UserUuid:      order.UserUUID.String(),
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: internal.ConvertPaymentMethodToProtobuf(req.PaymentMethod),
	})
	if err != nil {
		log.Printf("Payment service call failed for order %s: %v", orderUUIDStr, err)
		return internal.NewPayOrderInternalError(
			fmt.Sprintf("Payment failed: %v", err),
		), nil
	}

	transactionUUID, err := internal.ValidateTransactionUUID(paymentResp.TransactionUuid)
	if err != nil {
		log.Printf("Invalid transaction UUID received: %v", err)
		return internal.NewPayOrderInternalError(err.Error()), nil
	}

	order.TransactionUUID = &transactionUUID
	convertedPaymentMethod := internal.ConvertPaymentMethodToOrderDto(req.PaymentMethod)
	order.PaymentMethod = &convertedPaymentMethod
	order.Status = orderv1.OrderStatusPAID

	log.Printf("Order %s paid successfully, transaction %s", orderUUIDStr, transactionUUID.String())

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

/*
CancelOrder отменяет заказ.
Проверяет возможность отмены и обновляет статус заказа на CANCELLED.
*/
func (s *Server) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	orderUUIDStr := params.OrderUUID.String()
	log.Printf("Cancelling order %s", orderUUIDStr)

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderUUIDStr]
	if !exists {
		log.Printf("Order %s not found for cancellation", orderUUIDStr)
		return internal.NewOrderNotFoundError("Order not found"), nil
	}

	if err := internal.ValidateOrderCanBeCancelled(order); err != nil {
		log.Printf("Order %s cannot be cancelled: %v", orderUUIDStr, err)
		return internal.NewOrderConflictError(err.Error()), nil
	}

	order.Status = orderv1.OrderStatusCANCELLED
	log.Printf("Order %s cancelled successfully", orderUUIDStr)
	return &orderv1.CancelOrderNoContent{}, nil
}

// NewError создает стандартизированный ответ с ошибкой для неизвестных ошибок.
func (s *Server) NewError(ctx context.Context, err error) *orderv1.GenericErrorStatusCode {
	log.Printf("Unexpected error: %v", err)
	return internal.NewOrderGenericError(500, err.Error())
}

func main() {
	log.Println("Starting OrderService...")

	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	srv, err := orderv1.NewServer(server)
	if err != nil {
		log.Fatalf("Failed to create ogen server: %v", err)
	}

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           srv,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Println("OrderService listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve: %v", err)
	}
}
