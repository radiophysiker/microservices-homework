package encoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/radiophysiker/microservices-homework/assembly/internal/model"
	eventspb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/events/v1"
)

// EncodeShipAssembled сериализует событие ShipAssembled из доменной модели в protobuf байты
func EncodeShipAssembled(shipAssembled model.ShipAssembled) ([]byte, error) {
	pb := &eventspb.ShipAssembled{
		EventUuid:    shipAssembled.EventUUID.String(),
		OrderUuid:    shipAssembled.OrderUUID.String(),
		UserUuid:     shipAssembled.UserUUID.String(),
		BuildTimeSec: shipAssembled.BuildTimeSec,
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ShipAssembled: %w", err)
	}

	return data, nil
}
