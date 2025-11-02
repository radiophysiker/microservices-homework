package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiv1 "github.com/radiophysiker/microservices-homework/inventory/internal/api/inventory/v1"
	"github.com/radiophysiker/microservices-homework/inventory/internal/config"
	"github.com/radiophysiker/microservices-homework/inventory/internal/db"
	partRepo "github.com/radiophysiker/microservices-homework/inventory/internal/repository/part"
	partSvc "github.com/radiophysiker/microservices-homework/inventory/internal/service/part"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

func Run(ctx context.Context) error {
	mongoClient, collection, err := db.Connect(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Создаем зависимости
	partRepository := partRepo.NewRepository(collection)
	partService := partSvc.NewService(partRepository)
	api := apiv1.NewAPI(partService)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, api)
	health.RegisterService(grpcServer)
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("InventoryService gRPC server listening on", config.AppConfig().InventoryGRPC.Address())

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()

	log.Println("Server stopped")

	return nil
}
