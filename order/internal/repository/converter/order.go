package converter

import (
	"strings"

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

	serviceItems := make([]model.OrderItem, 0, len(repoOrder.Items))
	for _, it := range repoOrder.Items {
		serviceItems = append(serviceItems, model.OrderItem{
			PartUUID: it.PartUUID,
			Quantity: it.Quantity,
		})
	}

	return &model.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		Items:           serviceItems,
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

	repoItems := make([]repoModel.OrderItem, 0, len(serviceOrder.Items))
	for _, it := range serviceOrder.Items {
		repoItems = append(repoItems, repoModel.OrderItem{
			PartUUID: it.PartUUID,
			Quantity: it.Quantity,
		})
	}

	return &repoModel.Order{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		Items:           repoItems,
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

// StringToPaymentMethod конвертирует строку в PaymentMethod
func StringToPaymentMethod(s string) *repoModel.PaymentMethod {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "CARD":
		pm := repoModel.PaymentMethodCard
		return &pm
	case "SBP":
		pm := repoModel.PaymentMethodSBP
		return &pm
	case "CREDIT_CARD":
		pm := repoModel.PaymentMethodCreditCard
		return &pm
	case "INVESTOR_MONEY":
		pm := repoModel.PaymentMethodInvestorMoney
		return &pm
	default:
		pm := repoModel.PaymentMethodUnspecified
		return &pm
	}
}

// StringToOrderStatus конвертирует строку в Status
func StringToOrderStatus(s string) repoModel.Status {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "PENDING_PAYMENT":
		return repoModel.StatusPendingPayment
	case "PAID":
		return repoModel.StatusPaid
	case "CANCELLED":
		return repoModel.StatusCancelled
	default:
		return repoModel.StatusUnspecified
	}
}
