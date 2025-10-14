package service

import (
	"context"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// PartService представляет интерфейс для работы с деталями в сервисном слое
type PartService interface {
	// GetPart возвращает деталь по UUID
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	// ListParts возвращает список деталей с возможностью фильтрации
	ListParts(ctx context.Context, filter *model.Filter) ([]*model.Part, error)
}
