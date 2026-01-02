package kafka

import (
	"context"
	"fmt"
	"log"
	orderpb "order_service/proto/gen"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type consumer struct {
	reader      *kafka.Reader
	redisClient *redis.Client
}

type Consumer interface {
	StartResultListener(ctx context.Context)
}

func NewConsumer(reader *kafka.Reader) (Consumer, func() error) {
	c := &consumer{reader: reader}
	return c, c.reader.Close
}

func (c *consumer) StartResultListener(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		var result orderpb.OrderValidationResultEvent
		if err := proto.Unmarshal(m.Value, &result); err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			continue
		}
		status := "confirmed"
		if !result.IsValid {
			status = "Rejected"
		}
		msg := fmt.Sprintf("%s:%s", result.OrderId, status)
		c.redisClient.Publish(ctx, "order_updated", msg)
		log.Printf("Order %s processed. Status: %s", result.OrderId, status)
	}

}
