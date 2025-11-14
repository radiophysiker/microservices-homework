package encoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	eventspb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/events/v1"
)

func EncodeOrderPaid(orderPaid model.OrderPaid) ([]byte, error) {
	pb := &eventspb.OrderPaid{
		EventUuid:       orderPaid.EventUUID.String(),
		OrderUuid:       orderPaid.OrderUUID.String(),
		UserUuid:        orderPaid.UserUUID.String(),
		PaymentMethod:   converter.PaymentMethodToProtobuf(orderPaid.PaymentMethod),
		TransactionUuid: orderPaid.TransactionUUID.String(),
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OrderPaid: %w", err)
	}

	return data, nil
}
