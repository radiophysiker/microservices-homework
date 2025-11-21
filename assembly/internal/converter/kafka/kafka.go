package kafka

import "github.com/radiophysiker/microservices-homework/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (*model.OrderPaid, error)
}
