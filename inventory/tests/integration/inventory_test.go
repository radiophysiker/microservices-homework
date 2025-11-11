//go:build integration

package integration

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// TestGetPart тестирует получение детали по UUID через gRPC API
func (s *InventoryTestSuite) TestGetPart() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, cleanup, err := s.env.NewGRPCClient(ctx)
	s.Require().NoError(err, "Failed to create gRPC client")
	defer cleanup()

	testParts := GetTestParts()
	expectedPart := testParts[0]

	s.Run("success", func() {
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

		s.NotNil(part.Manufacturer, "Manufacturer should not be nil")
		s.Equal(expectedPart.Manufacturer.Name, part.Manufacturer.Name, "Manufacturer name should match")
		s.Equal(expectedPart.Manufacturer.Country, part.Manufacturer.Country, "Manufacturer country should match")
	})

	s.Run("not_found", func() {
		nonExistentReq := &pb.GetPartRequest{
			Uuid: "00000000-0000-0000-0000-000000000000",
		}

		_, err := client.GetPart(ctx, nonExistentReq)
		s.Error(err, "GetPart should return error for non-existent part")

		st, ok := status.FromError(err)
		s.True(ok, "error should be a gRPC status")
		s.Equal(codes.NotFound, st.Code(), "GetPart should return NotFound for non-existent part")
	})
}

// TestListParts тестирует получение списка деталей через gRPC API
func (s *InventoryTestSuite) TestListParts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, cleanup, err := s.env.NewGRPCClient(ctx)
	s.Require().NoError(err, "Failed to create gRPC client")
	defer cleanup()

	s.Run("all_parts", func() {
		req := &pb.ListPartsRequest{
			Filter: nil,
		}

		resp, err := client.ListParts(ctx, req)
		s.NoError(err, "ListParts should not return error")
		s.NotNil(resp, "Response should not be nil")
		s.NotNil(resp.Parts, "Parts should not be nil")
		s.GreaterOrEqual(len(resp.Parts), 4, "Should return at least 4 parts")
	})

	s.Run("by_category", func() {
		categoryFilter := &pb.ListPartsRequest{
			Filter: &pb.PartsFilter{
				Categories: []pb.Category{pb.Category_CATEGORY_ENGINE},
			},
		}

		resp, err := client.ListParts(ctx, categoryFilter)
		s.NoError(err, "ListParts with category filter should not return error")
		s.NotNil(resp.Parts, "Parts should not be nil")
		s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 engine part")

		for _, part := range resp.Parts {
			s.Equal(pb.Category_CATEGORY_ENGINE, part.Category, "All parts should be ENGINE category")
		}
	})

	s.Run("by_uuid", func() {
		testParts := GetTestParts()
		uuidFilter := &pb.ListPartsRequest{
			Filter: &pb.PartsFilter{
				Uuids: []string{testParts[0].UUID, testParts[1].UUID},
			},
		}

		resp, err := client.ListParts(ctx, uuidFilter)
		s.NoError(err, "ListParts with UUID filter should not return error")
		s.NotNil(resp.Parts, "Parts should not be nil")
		s.Equal(2, len(resp.Parts), "Should return exactly 2 parts")

		uuidMap := make(map[string]bool)
		for _, part := range resp.Parts {
			uuidMap[part.Uuid] = true
		}
		s.True(uuidMap[testParts[0].UUID], "Should contain first test part")
		s.True(uuidMap[testParts[1].UUID], "Should contain second test part")
	})

	s.Run("by_country", func() {
		countryFilter := &pb.ListPartsRequest{
			Filter: &pb.PartsFilter{
				ManufacturerCountries: []string{"USA"},
			},
		}

		resp, err := client.ListParts(ctx, countryFilter)
		s.NoError(err, "ListParts with country filter should not return error")
		s.NotNil(resp.Parts, "Parts should not be nil")
		s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 part from USA")

		for _, part := range resp.Parts {
			s.NotNil(part.Manufacturer, "Manufacturer should not be nil")
			s.Equal("USA", part.Manufacturer.Country, "All parts should be from USA")
		}
	})

	s.Run("by_tags", func() {
		tagFilter := &pb.ListPartsRequest{
			Filter: &pb.PartsFilter{
				Tags: []string{"engine"},
			},
		}

		resp, err := client.ListParts(ctx, tagFilter)
		s.NoError(err, "ListParts with tag filter should not return error")
		s.NotNil(resp.Parts, "Parts should not be nil")
		s.GreaterOrEqual(len(resp.Parts), 1, "Should return at least 1 part with 'engine' tag")

		for _, part := range resp.Parts {
			s.Contains(part.Tags, "engine", "All parts should have 'engine' tag")
		}
	})

	s.Run("empty_result", func() {
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
	})
}
