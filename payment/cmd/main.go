package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiv1 "github.com/radiophysiker/microservices-homework/payment/internal/api/payment/v1"
	paymentSvc "github.com/radiophysiker/microservices-homework/payment/internal/service/payment"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

func main() {
	paymentService := paymentSvc.NewService()
	api := apiv1.NewAPI(paymentService)

	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, api)
	reflection.Register(s)

	log.Println("PaymentService listening on :50052")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
