package kafka

import (
	"context"

	"github.com/olezhek28/microservices-course-olezhek-solution/platform/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
