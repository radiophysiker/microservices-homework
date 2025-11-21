package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/radiophysiker/microservices-homework/assembly/internal/model"
	eventspb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/events/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (*model.OrderPaid, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty message data")
	}

	var pb eventspb.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OrderPaid: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.GetEventUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid event_uuid: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.GetOrderUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid order_uuid: %w", err)
	}

	userUUID, err := uuid.Parse(pb.GetUserUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid user_uuid: %w", err)
	}

	transactionUUID, err := uuid.Parse(pb.GetTransactionUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_uuid: %w", err)
	}

	paymentMethod := paymentMethodFromProtobuf(pb.GetPaymentMethod())

	return &model.OrderPaid{
		EventUUID:       eventUUID,
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transactionUUID,
	}, nil
}

func paymentMethodFromProtobuf(pm paymentpb.PaymentMethod) model.PaymentMethod {
	switch pm {
	case paymentpb.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentpb.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentpb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentpb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
