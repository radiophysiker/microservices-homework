package decoder

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	eventspb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type Decoder struct{}

func NewOrderAssembledDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(data []byte) (*model.ShipAssembled, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty message data")
	}

	var pb eventspb.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ShipAssembled: %w", err)
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

	return &model.ShipAssembled{
		EventUUID:    eventUUID,
		OrderUUID:    orderUUID,
		UserUUID:     userUUID,
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
