package app

import (
	"context"
	"errors"
	"fmt"
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

	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	v1 "github.com/radiophysiker/microservices-homework/order/internal/api/order/v1"
	inventoryClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/payment/v1"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/order/internal/db"
	"github.com/radiophysiker/microservices-homework/order/internal/migrator"
	orderRepo "github.com/radiophysiker/microservices-homework/order/internal/repository/order"
	orderSvc "github.com/radiophysiker/microservices-homework/order/internal/service/order"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

func Run(ctx context.Context) error {
	// Подключаемся к базе данных
	pool, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Выполняем миграции
	if err = migrator.Run(ctx, pool, config.AppConfig().Migrations.Directory()); err != nil {
		return err
	}

	// External gRPC deps
	inventoryConn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.InventoryAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	paymentConn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.PaymentAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	// Создаем зависимости
	orderRepository := orderRepo.NewRepository(pool)
	inventoryClientInstance := inventoryClient.NewClient(inventorypb.NewInventoryServiceClient(inventoryConn))
	paymentClientInstance := paymentClient.NewClient(paymentpb.NewPaymentServiceClient(paymentConn))
	orderService := orderSvc.NewService(orderRepository, inventoryClientInstance, paymentClientInstance)

	// Создаем gRPC сервер для внутреннего использования
	grpcServer := grpc.NewServer()
	orderServiceServer := v1.NewAPI(orderService)
	orderpb.RegisterOrderServiceServer(grpcServer, orderServiceServer)
	health.RegisterService(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		grpcAddr := config.AppConfig().OrderGRPC.Address()

		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("OrderService gRPC server listening on ", grpcAddr)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Создаем gRPC Gateway для HTTP API
	ctx, cancel := context.WithCancel(ctx)

	mux := runtime.NewServeMux()

	// Регистрируем gRPC Gateway
	err = orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, config.AppConfig().OrderGRPC.Address(), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		cancel()
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// Запускаем HTTP сервер
	httpAddr := config.AppConfig().OrderHTTP.Address()
	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Println("OrderService HTTP Gateway listening on ", httpAddr)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	httpCtx, httpCancel := context.WithTimeout(ctx, 30*time.Second)
	defer httpCancel()

	if err := httpServer.Shutdown(httpCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()

	// Соединения уже закрыты через defer функции выше
	log.Println("Servers stopped")

	return nil
}
