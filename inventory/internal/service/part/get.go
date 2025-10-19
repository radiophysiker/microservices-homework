package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// GetPart возвращает деталь по UUID
func (s *Service) GetPart(ctx context.Context, id string) (*model.Part, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, model.NewInvalidUUIDError(id)
	}

	part, err := s.partRepository.GetPart(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return part, nil
}
