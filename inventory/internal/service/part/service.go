package part

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository"
)

// Service реализует интерфейс PartService
type Service struct {
	partRepository repository.PartRepository
}

// NewService создает новый экземпляр Service
func NewService(partRepository repository.PartRepository) *Service {
	return &Service{
		partRepository: partRepository,
	}
}
