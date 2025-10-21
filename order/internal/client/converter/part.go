package converter

import (
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// ToModelPart конвертирует protobuf Part в доменную модель Part
func ToModelPart(protoPart *inventorypb.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	return &model.Part{
		UUID:  protoPart.GetUuid(),
		Name:  protoPart.GetName(),
		Price: protoPart.GetPrice(),
	}
}

// ToModelParts конвертирует слайс protobuf Parts в слайс доменных моделей Parts
func ToModelParts(protoParts []*inventorypb.Part) []*model.Part {
	if len(protoParts) == 0 {
		return nil
	}

	parts := make([]*model.Part, len(protoParts))
	for i, protoPart := range protoParts {
		parts[i] = ToModelPart(protoPart)
	}

	return parts
}

// ToProtobufPartUUIDs конвертирует слайс UUID в слайс строк для protobuf
func ToProtobufPartUUIDs(partUUIDs []string) []string {
	return partUUIDs
}
