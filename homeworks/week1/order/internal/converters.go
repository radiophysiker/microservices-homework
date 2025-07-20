package internal

import (
	"github.com/google/uuid"
	orderv1 "github.com/radiophysiker/microservices-homework/week1/shared/pkg/openapi/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/week1/shared/pkg/proto/payment/v1"
)

// ConvertPaymentMethodToProtobuf конвертирует OpenAPI PaymentMethod в protobuf PaymentMethod
func ConvertPaymentMethodToProtobuf(openapi orderv1.PaymentMethod) paymentpb.PaymentMethod {
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

// ConvertPaymentMethodToOrderDto конвертирует PaymentMethod в OrderDtoPaymentMethod
func ConvertPaymentMethodToOrderDto(pm orderv1.PaymentMethod) orderv1.OrderDtoPaymentMethod {
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

// ConvertUUIDsToStrings конвертирует слайс UUID в слайс строк
func ConvertUUIDsToStrings(uuids []uuid.UUID) []string {
	strings := make([]string, len(uuids))
	for i, id := range uuids {
		strings[i] = id.String()
	}
	return strings
}

// ConvertOrderToDto конвертирует внутреннюю структуру Order в OrderDto
func ConvertOrderToDto(order *Order) *orderv1.OrderDto {
	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod *orderv1.NilOrderDtoPaymentMethod
	if order.PaymentMethod != nil {
		pm := orderv1.NewNilOrderDtoPaymentMethod(*order.PaymentMethod)
		paymentMethod = &pm
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          order.Status,
	}
}
