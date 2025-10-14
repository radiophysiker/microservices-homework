package converter

import (
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// ToOrderDto конвертирует модель service в DTO
func ToOrderDto(serviceOrder *model.Order) *orderv1.OrderDto {
	if serviceOrder == nil {
		return nil
	}

	var transactionUUID orderv1.OptNilUUID
	if serviceOrder.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*serviceOrder.TransactionUUID)
	}

	var paymentMethod *orderv1.NilOrderDtoPaymentMethod
	if serviceOrder.PaymentMethod != nil {
		pm := orderv1.NewNilOrderDtoPaymentMethod(*serviceOrder.PaymentMethod)
		paymentMethod = &pm
	}

	return &orderv1.OrderDto{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUuids:       serviceOrder.PartUUIDs,
		TotalPrice:      serviceOrder.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          serviceOrder.Status,
	}
}

// PaymentMethodToProtobuf конвертирует OpenAPI PaymentMethod в protobuf PaymentMethod
func PaymentMethodToProtobuf(openapi orderv1.PaymentMethod) paymentpb.PaymentMethod {
	switch openapi {
	case orderv1.PaymentMethodCARD:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentpb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// PaymentMethodToOrderDto конвертирует PaymentMethod в OrderDtoPaymentMethod
func PaymentMethodToOrderDto(pm orderv1.PaymentMethod) orderv1.OrderDtoPaymentMethod {
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
