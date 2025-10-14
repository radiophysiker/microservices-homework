package part

import (
	"context"
	"fmt"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// GetPart возвращает деталь по UUID
func (s *Service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	if uuid == "" {
		return nil, model.NewInvalidUUIDError(uuid)
	}

	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return part, nil
}
