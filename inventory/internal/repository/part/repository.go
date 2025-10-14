package part

import (
	"sync"

	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
)

// Repository реализует интерфейс PartRepository
type Repository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

// NewRepository создает новый экземпляр Repository
func NewRepository() *Repository {
	repo := &Repository{
		parts: make(map[string]*repoModel.Part),
	}
	repo.initTestData()
	return repo
}
