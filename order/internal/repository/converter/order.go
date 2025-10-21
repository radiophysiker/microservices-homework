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

	var paymentMethod *model.PaymentMethod

	if repoOrder.PaymentMethod != nil {
		pm := toServicePaymentMethod(*repoOrder.PaymentMethod)
		paymentMethod = &pm
	}

	return &model.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUUIDs:       repoOrder.PartUUIDs,
		TotalPrice:      repoOrder.TotalPrice,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          toServiceStatus(repoOrder.Status),
	}
}

// ToRepoOrder конвертирует модель service в модель repository
func ToRepoOrder(serviceOrder *model.Order) *repoModel.Order {
	if serviceOrder == nil {
		return nil
	}

	var paymentMethod *repoModel.PaymentMethod

	if serviceOrder.PaymentMethod != nil {
		pm := toRepoPaymentMethod(*serviceOrder.PaymentMethod)
		paymentMethod = &pm
	}

	return &repoModel.Order{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUUIDs:       serviceOrder.PartUUIDs,
		TotalPrice:      serviceOrder.TotalPrice,
		TransactionUUID: serviceOrder.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          toRepoStatus(serviceOrder.Status),
	}
}

func toServicePaymentMethod(pm repoModel.PaymentMethod) model.PaymentMethod {
	switch pm {
	case repoModel.PaymentMethodCard:
		return model.PaymentMethodCard
	case repoModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repoModel.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case repoModel.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func toRepoPaymentMethod(pm model.PaymentMethod) repoModel.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return repoModel.PaymentMethodCard
	case model.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return repoModel.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return repoModel.PaymentMethodInvestorMoney
	default:
		return repoModel.PaymentMethodUnspecified
	}
}

func toServiceStatus(s repoModel.Status) model.Status {
	switch s {
	case repoModel.StatusPendingPayment:
		return model.StatusPendingPayment
	case repoModel.StatusPaid:
		return model.StatusPaid
	case repoModel.StatusCancelled:
		return model.StatusCancelled
	default:
		return model.StatusUnspecified
	}
}

func toRepoStatus(s model.Status) repoModel.Status {
	switch s {
	case model.StatusPendingPayment:
		return repoModel.StatusPendingPayment
	case model.StatusPaid:
		return repoModel.StatusPaid
	case model.StatusCancelled:
		return repoModel.StatusCancelled
	default:
		return repoModel.StatusUnspecified
	}
}
