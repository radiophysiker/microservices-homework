package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// ToProtoPart converts domain model.Part to protobuf Part for the transport layer
func ToProtoPart(p *model.Part) *pb.Part {
	if p == nil {
		return nil
	}

	return &pb.Part{
		Uuid:         p.UUID,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price,
		Category:     toProtoCategory(p.Category),
		Dimensions:   toProtoDimensions(p.Dimensions),
		Manufacturer: toProtoManufacturer(p.Manufacturer),
		Tags:         p.Tags,
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}
}

// ToProtoParts конвертирует слайс моделей service в слайс protobuf моделей
func ToProtoParts(parts []*model.Part) []*pb.Part {
	if parts == nil {
		return []*pb.Part{}
	}

	protoParts := make([]*pb.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, ToProtoPart(part))
	}

	return protoParts
}

// ToModelFilter конвертирует protobuf фильтр в модель service
func ToModelFilter(filter *pb.PartsFilter) *model.Filter {
	if filter == nil {
		return nil
	}

	return &model.Filter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            toModelCategories(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func toProtoCategory(c model.Category) pb.Category {
	switch c {
	case model.CategoryEngine:
		return pb.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return pb.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return pb.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return pb.Category_CATEGORY_WING
	default:
		return pb.Category_CATEGORY_UNSPECIFIED
	}
}

func toModelCategory(c pb.Category) model.Category {
	switch c {
	case pb.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case pb.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case pb.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case pb.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func toModelCategories(categories []pb.Category) []model.Category {
	if categories == nil {
		return nil
	}

	result := make([]model.Category, 0, len(categories))
	for _, c := range categories {
		result = append(result, toModelCategory(c))
	}

	return result
}

func toProtoDimensions(d *model.Dimensions) *pb.Dimensions {
	if d == nil {
		return nil
	}

	return &pb.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func toProtoManufacturer(m *model.Manufacturer) *pb.Manufacturer {
	if m == nil {
		return nil
	}

	return &pb.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}
