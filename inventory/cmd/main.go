package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiv1 "github.com/radiophysiker/microservices-homework/inventory/internal/api/inventory/v1"
	partRepo "github.com/radiophysiker/microservices-homework/inventory/internal/repository/part"
	partSvc "github.com/radiophysiker/microservices-homework/inventory/internal/service/part"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

func main() {
	partRepository := partRepo.NewRepository()
	partService := partSvc.NewService(partRepository)
	api := apiv1.NewAPI(partService)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	log.Println("InventoryService listening on :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
