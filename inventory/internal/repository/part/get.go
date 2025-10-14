package part

import (
	"context"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
)

// GetPart возвращает деталь по UUID
func (r *Repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) { // ctx is currently unused; kept for interface compatibility
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.NewPartNotFoundError(uuid)
	}
	return converter.ToServicePart(part), nil
}
