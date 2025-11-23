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

func NewConsumer(reader *kafka.Reader) (Consumer, func() error) {
	c := &consumer{reader: reader}
	return c, c.reader.Close
}

func (p *consumer) Start(ctx context.Context, key, value []byte) error {
	for {
		msg, err := p.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		log.Println("message consumed")
		return nil
	}
}
