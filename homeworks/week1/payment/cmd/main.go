package main

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	pb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func newServer() *server {
	return &server{}
}

func (s *server) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	transactionUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &pb.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, newServer())
	reflection.Register(s)

	log.Println("PaymentService listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
