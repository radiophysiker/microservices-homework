package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	apiv1 "github.com/radiophysiker/microservices-homework/order/internal/api/order/v1"
	inventoryClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/radiophysiker/microservices-homework/order/internal/repository/order"
	orderSvc "github.com/radiophysiker/microservices-homework/order/internal/service/order"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}

	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	orderRepository := orderRepo.NewRepository()
	inventoryClientInstance := inventoryClient.NewClient(inventorypb.NewInventoryServiceClient(inventoryConn))
	paymentClientInstance := paymentClient.NewClient(paymentpb.NewPaymentServiceClient(paymentConn))
	orderService := orderSvc.NewService(orderRepository, inventoryClientInstance, paymentClientInstance)
	api := apiv1.NewAPI(orderService)

	srv, err := orderv1.NewServer(api)
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
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to serve: %v", err)
	}
}
