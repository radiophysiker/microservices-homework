package converter

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoPart converts domain model.Part to protobuf Part for the transport layer
func ToProtoPart(m *model.Part) *pb.Part {
	if m == nil {
		return nil
	}

	return &pb.Part{
		Uuid:          m.UUID,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: int64(m.StockQuantity),
		Category:      m.Category,
		Dimensions:    toProtoDimensions(m.Dimensions),
		Manufacturer:  toProtoManufacturer(m.Manufacturer),
		Tags:          m.Tags,
		// Metadata is omitted as it's not present in the domain model
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

// ToProtoParts конвертирует слайс моделей service в слайс protobuf моделей
func ToProtoParts(serviceParts []*model.Part) []*pb.Part {
	if serviceParts == nil {
		return nil
	}

	protoParts := make([]*pb.Part, 0, len(serviceParts))
	for _, servicePart := range serviceParts {
		protoParts = append(protoParts, ToProtoPart(servicePart))
	}

	return protoParts
}

// ToServiceFilter конвертирует protobuf фильтр в модель service
func ToServiceFilter(protoFilter *pb.PartsFilter) *model.Filter {
	if protoFilter == nil {
		return nil
	}

	return &model.Filter{
		UUIDs:                 protoFilter.Uuids,
		Names:                 protoFilter.Names,
		Categories:            protoFilter.Categories,
		ManufacturerCountries: protoFilter.ManufacturerCountries,
		Tags:                  protoFilter.Tags,
	}
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
