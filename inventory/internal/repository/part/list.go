package part

import (
	"context"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/radiophysiker/microservices-homework/inventory/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей с возможностью фильтрации
func (r *Repository) ListParts(_ context.Context, filter *model.Filter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Создаем копию всех частей
	allParts := make([]*repoModel.Part, 0, len(r.parts))
	for _, part := range r.parts {
		allParts = append(allParts, part)
	}

	// Если фильтр не задан, возвращаем все части
	if filter == nil || isEmptyFilter(filter) {
		return converter.ToServiceParts(allParts), nil
	}

	filteredParts := allParts

	if len(filter.UUIDs) > 0 {
		filteredParts = filterPartsByUUIDs(filteredParts, filter.UUIDs)
	}

	if len(filter.Names) > 0 {
		filteredParts = filterPartsByNames(filteredParts, filter.Names)
	}

	if len(filter.Categories) > 0 {
		filteredParts = filterPartsByCategories(filteredParts, filter.Categories)
	}

	if len(filter.ManufacturerCountries) > 0 {
		filteredParts = filterPartsByManufacturerCountries(filteredParts, filter.ManufacturerCountries)
	}

	if len(filter.Tags) > 0 {
		filteredParts = filterPartsByTags(filteredParts, filter.Tags)
	}

	return converter.ToServiceParts(filteredParts), nil
}

// isEmptyFilter проверяет, пуст ли фильтр
func isEmptyFilter(filter *model.Filter) bool {
	return len(filter.UUIDs) == 0 && len(filter.Names) == 0 &&
		len(filter.Categories) == 0 && len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

// filterPartsByUUIDs фильтрует детали по UUID
func filterPartsByUUIDs(parts []*repoModel.Part, uuids []string) []*repoModel.Part {
	filtered := make([]*repoModel.Part, 0, len(uuids))

	for _, part := range parts {
		for _, uuid := range uuids {
			if part.UUID == uuid {
				filtered = append(filtered, part)
				break
			}
		}
	}

	return filtered
}

// filterPartsByNames фильтрует детали по названиям
func filterPartsByNames(parts []*repoModel.Part, names []string) []*repoModel.Part {
	filtered := make([]*repoModel.Part, 0, len(names))

	for _, part := range parts {
		for _, name := range names {
			if part.Name == name {
				filtered = append(filtered, part)
				break
			}
		}
	}

	return filtered
}

// filterPartsByCategories фильтрует детали по категориям
func filterPartsByCategories(parts []*repoModel.Part, categories []pb.Category) []*repoModel.Part {
	filtered := make([]*repoModel.Part, 0, len(categories))

	for _, part := range parts {
		for _, category := range categories {
			if part.Category == category {
				filtered = append(filtered, part)
				break
			}
		}
	}

	return filtered
}

// filterPartsByManufacturerCountries фильтрует детали по странам производителей
func filterPartsByManufacturerCountries(parts []*repoModel.Part, countries []string) []*repoModel.Part {
	filtered := make([]*repoModel.Part, 0, len(countries))

	for _, part := range parts {
		if part.Manufacturer != nil {
			for _, country := range countries {
				if part.Manufacturer.Country == country {
					filtered = append(filtered, part)
					break
				}
			}
		}
	}

	return filtered
}

// filterPartsByTags фильтрует детали по тегам
func filterPartsByTags(parts []*repoModel.Part, tags []string) []*repoModel.Part {
	filtered := make([]*repoModel.Part, 0, len(parts))

	for _, part := range parts {
		for _, partTag := range part.Tags {
			for _, tag := range tags {
				if partTag == tag {
					filtered = append(filtered, part)
					goto nextPart
				}
			}
		}

		continue
	nextPart:
		filtered = append(filtered, part)
	}

	return filtered
}
