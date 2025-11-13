package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
}

type Consumer interface {
	Start(ctx context.Context, key, value []byte) error
}

func NewConsumer(reader *kafka.Reader) Consumer {
	return &consumer{
		reader: reader,
	}
}

func (p *consumer) Start(ctx context.Context, key, value []byte) error {
	msg, err := p.reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	log.Println("message consumed")
	return nil
}
