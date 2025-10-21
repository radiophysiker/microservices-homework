package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/radiophysiker/microservices-homework/order/internal/api/order/v1"
	inventoryClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/radiophysiker/microservices-homework/order/internal/repository/order"
	orderSvc "github.com/radiophysiker/microservices-homework/order/internal/service/order"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

func main() {
	// Подключаемся к внешним сервисам
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
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("failed to close inventory connection: %v", closeErr)
		}

		log.Fatalf("failed to connect to payment service: %v", err)
	}

	// Создаем зависимости
	orderRepository := orderRepo.NewRepository()
	inventoryClientInstance := inventoryClient.NewClient(inventorypb.NewInventoryServiceClient(inventoryConn))
	paymentClientInstance := paymentClient.NewClient(paymentpb.NewPaymentServiceClient(paymentConn))
	orderService := orderSvc.NewService(orderRepository, inventoryClientInstance, paymentClientInstance)

	// Создаем gRPC сервер для внутреннего использования
	grpcServer := grpc.NewServer()
	orderServiceServer := v1.NewAPI(orderService)
	orderpb.RegisterOrderServiceServer(grpcServer, orderServiceServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		lis, err := net.Listen("tcp", "localhost:50053")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("OrderService gRPC server listening on localhost:50053")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Создаем gRPC Gateway для HTTP API
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	mux := runtime.NewServeMux()

	// Регистрируем gRPC Gateway
	err = orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		cancel()
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Запускаем HTTP сервер
	httpServer := &http.Server{
		Addr:              "localhost:8080",
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Println("OrderService HTTP Gateway listening on localhost:8080")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Отменяем контекст
	cancel()

	// Останавливаем HTTP сервер
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer httpCancel()

	if err := httpServer.Shutdown(httpCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()

	// Закрываем соединения
	if err := inventoryConn.Close(); err != nil {
		log.Printf("failed to close inventory connection: %v", err)
	}

	if err := paymentConn.Close(); err != nil {
		log.Printf("failed to close payment connection: %v", err)
	}

	log.Println("Servers stopped")
}
