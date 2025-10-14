package model

import (
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь в сервисном слое order service
type Part struct {
	UUID  string
	Name  string
	Price float64
}

// ToServicePart конвертирует protobuf Part в модель сервисного слоя
func ToServicePart(pbPart *inventorypb.Part) *Part {
	if pbPart == nil {
		return nil
	}

	return &Part{
		UUID:  pbPart.Uuid,
		Name:  pbPart.Name,
		Price: pbPart.Price,
	}
}

// ToServiceParts конвертирует слайс protobuf Parts в слайс моделей сервисного слоя
func ToServiceParts(pbParts []*inventorypb.Part) []*Part {
	if pbParts == nil {
		return nil
	}

	parts := make([]*Part, 0, len(pbParts))
	for _, pbPart := range pbParts {
		parts = append(parts, ToServicePart(pbPart))
	}

	return parts
}
