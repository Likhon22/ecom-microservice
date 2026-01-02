package kafka

import (
	"context"
	"fmt"
	orderpb "order_service/proto/gen"

	"google.golang.org/protobuf/proto"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

type Producer interface {
	PublishOrderCreated(ctx context.Context, event *orderpb.OrderCreatedEvent) error
}

func NewProducer(writer *kafka.Writer) (Producer, func() error) {
	p := &producer{writer: writer}
	return p, p.writer.Close
}

func (p *producer) PublishOrderCreated(ctx context.Context, event *orderpb.OrderCreatedEvent) error {
	value, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderId),
		Value: value,
	})
}
