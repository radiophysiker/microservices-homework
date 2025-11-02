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

	apiv1 "github.com/radiophysiker/microservices-homework/payment/internal/api/payment/v1"
	"github.com/radiophysiker/microservices-homework/payment/internal/config"
	paymentSvc "github.com/radiophysiker/microservices-homework/payment/internal/service/payment"
	"github.com/radiophysiker/microservices-homework/platform/pkg/grpc/health"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

func Run(ctx context.Context) error {
	// Создаем зависимости
	paymentService := paymentSvc.NewService()
	api := apiv1.NewAPI(paymentService)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, api)
	health.RegisterService(grpcServer)
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		addr := config.AppConfig().PaymentGRPC.Address()

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("PaymentService gRPC server listening on", addr)

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
