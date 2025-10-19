package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// ConverterTestSuite тестовый набор для конвертера репозитория
type ConverterTestSuite struct {
	suite.Suite
}

// TestConverterSuite запускает тестовый набор
func TestConverterSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}

// TestToServicePart проверяет конвертацию из repoModel.Part в model.Part
func (s *ConverterTestSuite) TestToServicePart() {
	now := time.Now()

	tests := []struct {
		name     string
		input    *repoModel.Part
		expected *model.Part
	}{
		{
			name:     "nil_input",
			input:    nil,
			expected: nil,
		},
		{
			name: "full_part_with_all_fields",
			input: &repoModel.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Test Part",
				Description:   "Test Description",
				Price:         100.50,
				StockQuantity: 10,
				Category:      pb.Category_CATEGORY_ENGINE,
				Dimensions: &repoModel.Dimensions{
					Length: 100.0,
					Width:  50.0,
					Height: 30.0,
					Weight: 25.5,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    "Test Manufacturer",
					Country: "USA",
					Website: "https://test.com",
				},
				Tags:      []string{"tag1", "tag2"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			expected: &model.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Test Part",
				Description:   "Test Description",
				Price:         100.50,
				StockQuantity: 10,
				Category:      pb.Category_CATEGORY_ENGINE,
				Dimensions: &model.Dimensions{
					Length: 100.0,
					Width:  50.0,
					Height: 30.0,
					Weight: 25.5,
				},
				Manufacturer: &model.Manufacturer{
					Name:    "Test Manufacturer",
					Country: "USA",
					Website: "https://test.com",
				},
				Tags:      []string{"tag1", "tag2"},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "part_without_optional_fields",
			input: &repoModel.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Simple Part",
				Description:   "Simple Description",
				Price:         50.0,
				StockQuantity: 5,
				Category:      pb.Category_CATEGORY_WING,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expected: &model.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Simple Part",
				Description:   "Simple Description",
				Price:         50.0,
				StockQuantity: 5,
				Category:      pb.Category_CATEGORY_WING,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := ToServicePart(tt.input)
			s.Equal(tt.expected, result)
		})
	}
}

// TestToServiceParts проверяет конвертацию слайса частей из repo в service
func (s *ConverterTestSuite) TestToServiceParts() {
	now := time.Now()

	tests := []struct {
		name     string
		input    []*repoModel.Part
		expected []*model.Part
	}{
		{
			name:     "nil_slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty_slice",
			input:    []*repoModel.Part{},
			expected: []*model.Part{},
		},
		{
			name: "slice_with_multiple_parts",
			input: []*repoModel.Part{
				{
					UUID:          "550e8400-e29b-41d4-a716-446655440001",
					Name:          "Part 1",
					Description:   "Description 1",
					Price:         100.0,
					StockQuantity: 10,
					Category:      pb.Category_CATEGORY_ENGINE,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
				{
					UUID:          "550e8400-e29b-41d4-a716-446655440002",
					Name:          "Part 2",
					Description:   "Description 2",
					Price:         200.0,
					StockQuantity: 20,
					Category:      pb.Category_CATEGORY_FUEL,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			expected: []*model.Part{
				{
					UUID:          "550e8400-e29b-41d4-a716-446655440001",
					Name:          "Part 1",
					Description:   "Description 1",
					Price:         100.0,
					StockQuantity: 10,
					Category:      pb.Category_CATEGORY_ENGINE,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
				{
					UUID:          "550e8400-e29b-41d4-a716-446655440002",
					Name:          "Part 2",
					Description:   "Description 2",
					Price:         200.0,
					StockQuantity: 20,
					Category:      pb.Category_CATEGORY_FUEL,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := ToServiceParts(tt.input)
			s.Equal(tt.expected, result)
		})
	}
}

// TestToRepoPart проверяет конвертацию из model.Part в repoModel.Part
func (s *ConverterTestSuite) TestToRepoPart() {
	now := time.Now()

	tests := []struct {
		name     string
		input    *model.Part
		expected *repoModel.Part
	}{
		{
			name:     "nil_input",
			input:    nil,
			expected: nil,
		},
		{
			name: "full_part_with_all_fields",
			input: &model.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Test Part",
				Description:   "Test Description",
				Price:         100.50,
				StockQuantity: 10,
				Category:      pb.Category_CATEGORY_ENGINE,
				Dimensions: &model.Dimensions{
					Length: 100.0,
					Width:  50.0,
					Height: 30.0,
					Weight: 25.5,
				},
				Manufacturer: &model.Manufacturer{
					Name:    "Test Manufacturer",
					Country: "USA",
					Website: "https://test.com",
				},
				Tags:      []string{"tag1", "tag2"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			expected: &repoModel.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Test Part",
				Description:   "Test Description",
				Price:         100.50,
				StockQuantity: 10,
				Category:      pb.Category_CATEGORY_ENGINE,
				Dimensions: &repoModel.Dimensions{
					Length: 100.0,
					Width:  50.0,
					Height: 30.0,
					Weight: 25.5,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    "Test Manufacturer",
					Country: "USA",
					Website: "https://test.com",
				},
				Tags:      []string{"tag1", "tag2"},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "part_without_optional_fields",
			input: &model.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Simple Part",
				Description:   "Simple Description",
				Price:         50.0,
				StockQuantity: 5,
				Category:      pb.Category_CATEGORY_WING,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expected: &repoModel.Part{
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Simple Part",
				Description:   "Simple Description",
				Price:         50.0,
				StockQuantity: 5,
				Category:      pb.Category_CATEGORY_WING,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := ToRepoPart(tt.input)
			s.Equal(tt.expected, result)
		})
	}
}
