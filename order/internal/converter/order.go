package converter

import (
	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
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

	partUUIDs := make([]uuid.UUID, 0, len(serviceOrder.Items))
	for _, it := range serviceOrder.Items {
		partUUIDs = append(partUUIDs, it.PartUUID)
	}

	return &orderv1.OrderDto{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUuids:       partUUIDs,
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

// PaymentMethodFromProtobuf конвертирует protobuf PaymentMethod в model.PaymentMethod
func PaymentMethodFromProtobuf(pm orderpb.PaymentMethod) model.PaymentMethod {
	switch pm {
	case orderpb.PaymentMethod_CARD:
		return model.PaymentMethodCard
	case orderpb.PaymentMethod_SBP:
		return model.PaymentMethodSBP
	case orderpb.PaymentMethod_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case orderpb.PaymentMethod_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

// PaymentMethodToOrderProtobuf конвертирует model.PaymentMethod в orderpb.PaymentMethod
func PaymentMethodToOrderProtobuf(pm model.PaymentMethod) orderpb.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return orderpb.PaymentMethod_CARD
	case model.PaymentMethodSBP:
		return orderpb.PaymentMethod_SBP
	case model.PaymentMethodCreditCard:
		return orderpb.PaymentMethod_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return orderpb.PaymentMethod_INVESTOR_MONEY
	default:
		return orderpb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// StatusToProtobuf конвертирует model.Status в protobuf OrderStatus
func StatusToProtobuf(s model.Status) orderpb.OrderStatus {
	switch s {
	case model.StatusPendingPayment:
		return orderpb.OrderStatus_ORDER_STATUS_PENDING_PAYMENT
	case model.StatusPaid:
		return orderpb.OrderStatus_ORDER_STATUS_PAID
	case model.StatusCancelled:
		return orderpb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// ToProtoOrder конвертирует доменную модель в protobuf GetOrderResponse
func ToProtoOrder(serviceOrder *model.Order) *orderpb.GetOrderResponse {
	if serviceOrder == nil {
		return nil
	}

	var transactionUUID *string

	if serviceOrder.TransactionUUID != nil {
		transactionUUIDStr := serviceOrder.TransactionUUID.String()
		transactionUUID = &transactionUUIDStr
	}

	var paymentMethod *orderpb.PaymentMethod

	if serviceOrder.PaymentMethod != nil {
		pm := PaymentMethodToOrderProtobuf(*serviceOrder.PaymentMethod)
		paymentMethod = &pm
	}

	partUUIDStrings := make([]string, 0, len(serviceOrder.Items))
	for _, it := range serviceOrder.Items {
		partUUIDStrings = append(partUUIDStrings, it.PartUUID.String())
	}

	return &orderpb.GetOrderResponse{
		OrderUuid:       serviceOrder.OrderUUID.String(),
		UserUuid:        serviceOrder.UserUUID.String(),
		PartUuids:       partUUIDStrings,
		TotalPrice:      serviceOrder.TotalPrice,
		TransactionUuid: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          StatusToProtobuf(serviceOrder.Status),
	}
}
