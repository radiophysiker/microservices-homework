package part

import (
	"context"
	"fmt"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// ListParts возвращает список деталей с возможностью фильтрации
func (s *Service) ListParts(ctx context.Context, filter *model.Filter) ([]*model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %w", err)
	}

	return parts, nil
}
