package kafka

import (
	"github.com/radiophysiker/microservices-homework/notification/internal/model"
)

// OrderPaidDecoder декодирует сообщения OrderPaid из Kafka
type OrderPaidDecoder interface {
	Decode(data []byte) (*model.OrderPaid, error)
}

// OrderAssembledDecoder декодирует сообщения ShipAssembled из Kafka
type OrderAssembledDecoder interface {
	Decode(data []byte) (*model.ShipAssembled, error)
}
