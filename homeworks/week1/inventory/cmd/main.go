package main

import (
	"context"
	"log"
	"net"
	"slices"
	"sync"

	pb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*pb.Part
}

func newServer() *server {
	s := &server{
		parts: make(map[string]*pb.Part),
	}
	s.initTestData()
	return s
}

func (s *server) initTestData() {
	now := timestamppb.Now()

	testParts := []*pb.Part{
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Главный двигатель V8",
			Description:   "Мощный ракетный двигатель для основной тяги",
			Price:         50000.00,
			StockQuantity: 10,
			Category:      pb.Category_CATEGORY_ENGINE,
			Dimensions: &pb.Dimensions{
				Length: 300.0,
				Width:  150.0,
				Height: 200.0,
				Weight: 5000.0,
			},
			Manufacturer: &pb.Manufacturer{
				Name:    "SpaceX Engines",
				Country: "USA",
				Website: "https://spacex.com",
			},
			Tags:      []string{"main", "powerful", "v8"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Топливный бак",
			Description:   "Герметичный топливный бак для ракетного топлива",
			Price:         15000.00,
			StockQuantity: 25,
			Category:      pb.Category_CATEGORY_FUEL,
			Dimensions: &pb.Dimensions{
				Length: 400.0,
				Width:  200.0,
				Height: 250.0,
				Weight: 1000.0,
			},
			Manufacturer: &pb.Manufacturer{
				Name:    "FuelTech GmbH",
				Country: "Germany",
				Website: "https://fueltech.de",
			},
			Tags:      []string{"fuel", "storage", "sealed"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Обзорный иллюминатор",
			Description:   "Прочный иллюминатор из закаленного стекла",
			Price:         3000.00,
			StockQuantity: 50,
			Category:      pb.Category_CATEGORY_PORTHOLE,
			Dimensions: &pb.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: &pb.Manufacturer{
				Name:    "ClearView Ltd",
				Country: "Japan",
				Website: "https://clearview.jp",
			},
			Tags:      []string{"view", "glass", "durable"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Стабилизирующее крыло",
			Description:   "Аэродинамическое крыло для стабилизации полета",
			Price:         8000.00,
			StockQuantity: 20,
			Category:      pb.Category_CATEGORY_WING,
			Dimensions: &pb.Dimensions{
				Length: 500.0,
				Width:  100.0,
				Height: 50.0,
				Weight: 800.0,
			},
			Manufacturer: &pb.Manufacturer{
				Name:    "AeroWings Corp",
				Country: "France",
				Website: "https://aerowings.fr",
			},
			Tags:      []string{"wing", "stabilizer", "aerodynamic"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, part := range testParts {
		s.parts[part.Uuid] = part
	}
}

func (s *server) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, exists := s.parts[req.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
	}

	return &pb.GetPartResponse{Part: part}, nil
}

func (s *server) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создаем копию всех частей
	allParts := make([]*pb.Part, 0, len(s.parts))
	for _, part := range s.parts {
		allParts = append(allParts, part)
	}

	// Если фильтр не задан, возвращаем все части
	filter := req.Filter
	if filter == nil || (len(filter.Uuids) == 0 && len(filter.Names) == 0 &&
		len(filter.Categories) == 0 && len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0) {
		return &pb.ListPartsResponse{Parts: allParts}, nil
	}

	filteredParts := allParts

	if len(filter.Uuids) > 0 {
		filteredParts = filterPartsByUuids(filteredParts, filter.Uuids)
	}

	if len(filter.Names) > 0 {
		filteredParts = filterPartsByName(filteredParts, filter.Names)
	}

	if len(filter.Categories) > 0 {
		filteredParts = filterPartsByCategory(filteredParts, filter.Categories)
	}

	if len(filter.ManufacturerCountries) > 0 {
		filteredParts = filterPartsByManufacturerCountry(filteredParts, filter.ManufacturerCountries)
	}

	if len(filter.Tags) > 0 {
		filteredParts = filterPartsByTag(filteredParts, filter.Tags)
	}

	return &pb.ListPartsResponse{Parts: filteredParts}, nil
}

func filterPartsByTag(filteredParts []*pb.Part, tags []string) []*pb.Part {
	filtered := make([]*pb.Part, 0, len(filteredParts))

	for _, part := range filteredParts {
		for _, partTag := range part.Tags {
			if slices.Contains(tags, partTag) {
				filtered = append(filtered, part)
				break
			}
		}
	}
	return filtered
}

func filterPartsByManufacturerCountry(filteredParts []*pb.Part, countries []string) []*pb.Part {
	filtered := make([]*pb.Part, 0, len(filteredParts))
	for _, part := range filteredParts {
		if part.Manufacturer != nil && slices.Contains(countries, part.Manufacturer.Country) {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func filterPartsByCategory(filteredParts []*pb.Part, categories []pb.Category) []*pb.Part {
	filtered := make([]*pb.Part, 0, len(filteredParts))
	for _, part := range filteredParts {
		if slices.Contains(categories, part.Category) {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func filterPartsByName(filteredParts []*pb.Part, s []string) []*pb.Part {
	filtered := make([]*pb.Part, 0, len(s))
	for _, part := range filteredParts {
		if slices.Contains(s, part.Name) {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func filterPartsByUuids(filteredParts []*pb.Part, s []string) []*pb.Part {
	filtered := make([]*pb.Part, 0, len(s))
	for _, part := range filteredParts {
		if slices.Contains(s, part.Uuid) {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, newServer())
	reflection.Register(s)

	log.Println("InventoryService listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
