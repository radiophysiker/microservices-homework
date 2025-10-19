package converter

import (
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// ToOpenAPIOrder конвертирует доменную модель в OpenAPI DTO
func ToOpenAPIOrder(serviceOrder *model.Order) *orderv1.OrderDto {
	if serviceOrder == nil {
		return nil
	}

	var transactionUUID orderv1.OptNilUUID
	if serviceOrder.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*serviceOrder.TransactionUUID)
	}

	var paymentMethod *orderv1.NilOrderDtoPaymentMethod

	if serviceOrder.PaymentMethod != nil {
		pm := toOpenAPIPaymentMethod(*serviceOrder.PaymentMethod)
		nilPm := orderv1.NewNilOrderDtoPaymentMethod(pm)
		paymentMethod = &nilPm
	}

	return &orderv1.OrderDto{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUuids:       serviceOrder.PartUUIDs,
		TotalPrice:      serviceOrder.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          toOpenAPIStatus(serviceOrder.Status),
	}
}

// toOpenAPIPaymentMethod конвертирует model.PaymentMethod в orderv1.OrderDtoPaymentMethod
func toOpenAPIPaymentMethod(pm model.PaymentMethod) orderv1.OrderDtoPaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return orderv1.OrderDtoPaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderv1.OrderDtoPaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderv1.OrderDtoPaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderv1.OrderDtoPaymentMethodINVESTORMONEY
	default:
		return orderv1.OrderDtoPaymentMethodUNKNOWN
	}
}

// ToModelPaymentMethod конвертирует orderv1.PaymentMethod в model.PaymentMethod
func ToModelPaymentMethod(pm orderv1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

// toOpenAPIStatus конвертирует model.Status в orderv1.OrderStatus
func toOpenAPIStatus(s model.Status) orderv1.OrderStatus {
	switch s {
	case model.StatusPendingPayment:
		return orderv1.OrderStatusPENDINGPAYMENT
	case model.StatusPaid:
		return orderv1.OrderStatusPAID
	case model.StatusCancelled:
		return orderv1.OrderStatusCANCELLED
	default:
		return orderv1.OrderStatusPENDINGPAYMENT
	}
}

// PaymentMethodToProtobuf конвертирует model.PaymentMethod в protobuf PaymentMethod
func PaymentMethodToProtobuf(pm model.PaymentMethod) paymentpb.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// PaymentMethodToOpenAPI конвертирует PaymentMethod в OpenAPI OrderDtoPaymentMethod
func PaymentMethodToOpenAPI(pm orderv1.PaymentMethod) orderv1.OrderDtoPaymentMethod {
	switch pm {
	case orderv1.PaymentMethodCARD:
		return orderv1.OrderDtoPaymentMethodCARD
	case orderv1.PaymentMethodSBP:
		return orderv1.OrderDtoPaymentMethodSBP
	case orderv1.PaymentMethodCREDITCARD:
		return orderv1.OrderDtoPaymentMethodCREDITCARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return orderv1.OrderDtoPaymentMethodINVESTORMONEY
	default:
		return orderv1.OrderDtoPaymentMethodUNKNOWN
	}
}
