package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

type Producer interface {
	Publish(ctx context.Context, key, value []byte) error
}

func NewProducer(writer *kafka.Writer) (Producer, func() error) {
	p := &producer{writer: writer}
	return p, p.writer.Close
}

func (p *producer) Publish(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}
	log.Println("message published")
	return nil
}
