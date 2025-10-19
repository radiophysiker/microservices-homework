package converter

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	repoModel "github.com/radiophysiker/microservices-homework/inventory/internal/repository/model"
)

// ToServicePart конвертирует модель repository в модель service
func ToServicePart(repoPart *repoModel.Part) *model.Part {
	if repoPart == nil {
		return nil
	}

	return &model.Part{
		UUID:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      toServiceCategory(repoPart.Category),
		Dimensions:    toServiceDimensions(repoPart.Dimensions),
		Manufacturer:  toServiceManufacturer(repoPart.Manufacturer),
		Tags:          repoPart.Tags,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}
}

// ToServiceParts конвертирует слайс моделей repository в слайс моделей service
func ToServiceParts(repoParts []*repoModel.Part) []*model.Part {
	if repoParts == nil {
		return nil
	}

	serviceParts := make([]*model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		serviceParts = append(serviceParts, ToServicePart(repoPart))
	}

	return serviceParts
}

// ToRepoPart конвертирует модель service в модель repository
func ToRepoPart(servicePart *model.Part) *repoModel.Part {
	if servicePart == nil {
		return nil
	}

	return &repoModel.Part{
		UUID:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      ToRepoCategory(servicePart.Category),
		Dimensions:    toRepoDimensions(servicePart.Dimensions),
		Manufacturer:  toRepoManufacturer(servicePart.Manufacturer),
		Tags:          servicePart.Tags,
		CreatedAt:     servicePart.CreatedAt,
		UpdatedAt:     servicePart.UpdatedAt,
	}
}

func toServiceCategory(c repoModel.Category) model.Category {
	switch c {
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

// ToRepoCategory конвертирует категорию из service модели в repository модель
func ToRepoCategory(c model.Category) repoModel.Category {
	switch c {
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnspecified
	}
}

func toServiceDimensions(repoDimensions *repoModel.Dimensions) *model.Dimensions {
	if repoDimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: repoDimensions.Length,
		Width:  repoDimensions.Width,
		Height: repoDimensions.Height,
		Weight: repoDimensions.Weight,
	}
}

func toServiceManufacturer(repoManufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if repoManufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    repoManufacturer.Name,
		Country: repoManufacturer.Country,
		Website: repoManufacturer.Website,
	}
}

func toRepoDimensions(serviceDimensions *model.Dimensions) *repoModel.Dimensions {
	if serviceDimensions == nil {
		return nil
	}

	return &repoModel.Dimensions{
		Length: serviceDimensions.Length,
		Width:  serviceDimensions.Width,
		Height: serviceDimensions.Height,
		Weight: serviceDimensions.Weight,
	}
}

func toRepoManufacturer(serviceManufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if serviceManufacturer == nil {
		return nil
	}

	return &repoModel.Manufacturer{
		Name:    serviceManufacturer.Name,
		Country: serviceManufacturer.Country,
		Website: serviceManufacturer.Website,
	}
}
