package converter

import (
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	repoModel "github.com/radiophysiker/microservices-homework/order/internal/repository/model"
)

// ToServiceOrder конвертирует модель repository в модель service
func ToServiceOrder(repoOrder *repoModel.Order) *model.Order {
	if repoOrder == nil {
		return nil
	}

	return &model.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUUIDs:       repoOrder.PartUUIDs,
		TotalPrice:      repoOrder.TotalPrice,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   repoOrder.PaymentMethod,
		Status:          repoOrder.Status,
	}
}

// ToRepoOrder конвертирует модель service в модель repository
func ToRepoOrder(serviceOrder *model.Order) *repoModel.Order {
	if serviceOrder == nil {
		return nil
	}

	return &repoModel.Order{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUUIDs:       serviceOrder.PartUUIDs,
		TotalPrice:      serviceOrder.TotalPrice,
		TransactionUUID: serviceOrder.TransactionUUID,
		PaymentMethod:   serviceOrder.PaymentMethod,
		Status:          serviceOrder.Status,
	}
}
