package part

import (
	"time"

	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// initTestData инициализирует тестовые данные
func (r *Repository) initTestData() {
	now := time.Now()

	testParts := []*repoModel.Part{
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Главный двигатель V8",
			Description:   "Мощный ракетный двигатель для основной тяги",
			Price:         50000.00,
			StockQuantity: 10,
			Category:      pb.Category_CATEGORY_ENGINE,
			Dimensions: &repoModel.Dimensions{
				Length: 300.0,
				Width:  150.0,
				Height: 200.0,
				Weight: 5000.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX Engines",
				Country: "USA",
				Website: "https://spacex.com",
			},
			Tags:      []string{"main", "powerful", "v8"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Топливный бак",
			Description:   "Герметичный топливный бак для ракетного топлива",
			Price:         15000.00,
			StockQuantity: 25,
			Category:      pb.Category_CATEGORY_FUEL,
			Dimensions: &repoModel.Dimensions{
				Length: 400.0,
				Width:  200.0,
				Height: 250.0,
				Weight: 1000.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "FuelTech GmbH",
				Country: "Germany",
				Website: "https://fueltech.de",
			},
			Tags:      []string{"fuel", "storage", "sealed"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Обзорный иллюминатор",
			Description:   "Прочный иллюминатор из закаленного стекла",
			Price:         3000.00,
			StockQuantity: 50,
			Category:      pb.Category_CATEGORY_PORTHOLE,
			Dimensions: &repoModel.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ClearView Ltd",
				Country: "Japan",
				Website: "https://clearview.jp",
			},
			Tags:      []string{"view", "glass", "durable"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Стабилизирующее крыло",
			Description:   "Аэродинамическое крыло для стабилизации полета",
			Price:         8000.00,
			StockQuantity: 20,
			Category:      pb.Category_CATEGORY_WING,
			Dimensions: &repoModel.Dimensions{
				Length: 500.0,
				Width:  100.0,
				Height: 50.0,
				Weight: 800.0,
			},
			Manufacturer: &repoModel.Manufacturer{
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
		r.parts[part.UUID] = part
	}
}
