//go:build integration

package integration

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// TestGetPart тестирует получение детали по UUID через gRPC API
func (s *InventoryTestSuite) TestGetPart() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient(
		s.env.AppAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err, "Failed to connect to gRPC server")
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewInventoryServiceClient(conn)

	// Тест 1: Получение существующей детали
	testParts := GetTestParts()
	expectedPart := testParts[0] // Main Engine V8

	req := &pb.GetPartRequest{
		Uuid: expectedPart.UUID,
	}

	resp, err := client.GetPart(ctx, req)
	s.NoError(err, "GetPart should not return error")
	s.NotNil(resp, "Response should not be nil")
	s.NotNil(resp.Part, "Part should not be nil")

	part := resp.Part
	s.Equal(expectedPart.UUID, part.Uuid, "UUID should match")
	s.Equal(expectedPart.Name, part.Name, "Name should match")
	s.Equal(expectedPart.Description, part.Description, "Description should match")
	s.Equal(expectedPart.Price, part.Price, "Price should match")
	s.Equal(pb.Category_CATEGORY_ENGINE, part.Category, "Category should match")
	s.Equal(expectedPart.Tags, part.Tags, "Tags should match")

	// Проверяем производителя
	s.NotNil(part.Manufacturer, "Manufacturer should not be nil")
	s.Equal(expectedPart.Manufacturer.Name, part.Manufacturer.Name, "Manufacturer name should match")
	s.Equal(expectedPart.Manufacturer.Country, part.Manufacturer.Country, "Manufacturer country should match")

	// Тест 2: Получение несуществующей детали
	nonExistentReq := &pb.GetPartRequest{
		Uuid: "00000000-0000-0000-0000-000000000000",
	}

	_, err = client.GetPart(ctx, nonExistentReq)
	s.Error(err, "GetPart should return error for non-existent part")
}

// TestListParts тестирует получение списка деталей через gRPC API
func (s *InventoryTestSuite) TestListParts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient(
		s.env.AppAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err, "Failed to connect to gRPC server")
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewInventoryServiceClient(conn)

	// Тест 1: Получение всех деталей без фильтра
	req := &pb.ListPartsRequest{
		Filter: nil,
	}

	resp, err := client.ListParts(ctx, req)
	s.NoError(err, "ListParts should not return error")
	s.NotNil(resp, "Response should not be nil")
	s.NotNil(resp.Parts, "Parts should not be nil")
	s.GreaterOrEqual(len(resp.Parts), 4, "Should return at least 4 parts")

	// Тест 2: Фильтрация по категории
	categoryFilter := &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			Categories: []pb.Category{pb.Category_CATEGORY_ENGINE},
		},
	}

	resp, err = client.ListParts(ctx, categoryFilter)
	s.NoError(err, "ListParts with category filter should not return error")
	s.NotNil(resp.Parts, "Parts should not be nil")
	s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 engine part")

	// Проверяем, что все возвращенные детали имеют категорию ENGINE
	for _, part := range resp.Parts {
		s.Equal(pb.Category_CATEGORY_ENGINE, part.Category, "All parts should be ENGINE category")
	}

	// Тест 3: Фильтрация по UUID
	testParts := GetTestParts()
	uuidFilter := &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			Uuids: []string{testParts[0].UUID, testParts[1].UUID},
		},
	}

	resp, err = client.ListParts(ctx, uuidFilter)
	s.NoError(err, "ListParts with UUID filter should not return error")
	s.NotNil(resp.Parts, "Parts should not be nil")
	s.Equal(2, len(resp.Parts), "Should return exactly 2 parts")

	// Проверяем, что возвращены правильные детали
	uuidMap := make(map[string]bool)
	for _, part := range resp.Parts {
		uuidMap[part.Uuid] = true
	}
	s.True(uuidMap[testParts[0].UUID], "Should contain first test part")
	s.True(uuidMap[testParts[1].UUID], "Should contain second test part")

	// Тест 4: Фильтрация по стране производителя
	countryFilter := &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			ManufacturerCountries: []string{"USA"},
		},
	}

	resp, err = client.ListParts(ctx, countryFilter)
	s.NoError(err, "ListParts with country filter should not return error")
	s.NotNil(resp.Parts, "Parts should not be nil")
	s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 part from USA")

	// Проверяем, что все возвращенные детали из USA
	for _, part := range resp.Parts {
		s.NotNil(part.Manufacturer, "Manufacturer should not be nil")
		s.Equal("USA", part.Manufacturer.Country, "All parts should be from USA")
	}

	// Тест 5: Фильтрация по тегам
	tagFilter := &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			Tags: []string{"engine"},
		},
	}

	resp, err = client.ListParts(ctx, tagFilter)
	s.NoError(err, "ListParts with tag filter should not return error")
	s.NotNil(resp.Parts, "Parts should not be nil")
	s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 part with 'engine' tag")

	// Проверяем, что все возвращенные детали содержат тег "engine"
	for _, part := range resp.Parts {
		s.Contains(part.Tags, "engine", "All parts should have 'engine' tag")
	}
}

// TestListPartsEmptyResult тестирует фильтрацию, которая не возвращает результатов
func (s *InventoryTestSuite) TestListPartsEmptyResult() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient(
		s.env.AppAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err, "Failed to connect to gRPC server")
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewInventoryServiceClient(conn)

	// Фильтр, который не должен вернуть результатов
	req := &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			Uuids: []string{"99999999-9999-9999-9999-999999999999"},
		},
	}

	resp, err := client.ListParts(ctx, req)
	s.NoError(err, "ListParts should not return error even for empty result")
	s.NotNil(resp, "Response should not be nil")
	parts := resp.GetParts()
	if parts == nil {
		parts = []*pb.Part{}
	}
	s.Equal(0, len(parts), "Should return empty list for non-existent UUID")
}
