package part

import (
	"context"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
)

// GetPart возвращает деталь по UUID
func (r *Repository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.NewErrPartNotFound(uuid)
	}

	return converter.ToServicePart(part), nil
}
