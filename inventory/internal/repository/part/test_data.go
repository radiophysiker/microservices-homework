package part

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// Тестовые UUID для деталей
const (
	mainEngineV8UUID = "123e4567-e89b-12d3-a456-426614174000"
	fuelTankUUID     = "223e4567-e89b-12d3-a456-426614174001"
	wingUUID         = "323e4567-e89b-12d3-a456-426614174002"
	cockpitUUID      = "423e4567-e89b-12d3-a456-426614174003"
)

// testParts содержит тестовые данные деталей
var testParts = []*model.Part{
	{
		UUID:        mainEngineV8UUID,
		Name:        "Main Engine V8",
		Description: "High-performance rocket engine",
		Price:       50000,
		Category:    model.CategoryEngine,
		Manufacturer: &model.Manufacturer{
			Name:    "SpaceTech",
			Country: "USA",
		},
		Tags: []string{"engine", "propulsion", "v8"},
	},
	{
		UUID:        fuelTankUUID,
		Name:        "Fuel Tank",
		Description: "Large capacity fuel storage",
		Price:       15000,
		Category:    model.CategoryFuel,
		Manufacturer: &model.Manufacturer{
			Name:    "FuelCorp",
			Country: "Germany",
		},
		Tags: []string{"fuel", "storage", "tank"},
	},
	{
		UUID:        wingUUID,
		Name:        "Wing Assembly",
		Description: "Aerodynamic wing structure",
		Price:       25000,
		Category:    model.CategoryWing,
		Manufacturer: &model.Manufacturer{
			Name:    "AeroParts",
			Country: "France",
		},
		Tags: []string{"wing", "structure", "aerodynamics"},
	},
	{
		UUID:        cockpitUUID,
		Name:        "Cockpit Module",
		Description: "Pilot control center",
		Price:       35000,
		Category:    model.CategoryPorthole,
		Manufacturer: &model.Manufacturer{
			Name:    "ControlTech",
			Country: "Japan",
		},
		Tags: []string{"cockpit", "control", "pilot"},
	},
}

// GetTestParts возвращает тестовые данные деталей для использования в тестах
func GetTestParts() []*model.Part {
	return testParts
}
