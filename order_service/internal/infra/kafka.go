package infra

import (
	"github.com/segmentio/kafka-go"
)

type KafkaInfra struct {
	Brokers []string
}

func NewKafkaInfra(brokers []string) *KafkaInfra {

	return &KafkaInfra{
		Brokers: brokers,
	}
}

func (k *KafkaInfra) Writer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.Brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *KafkaInfra) Reader(topic, groupId string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		GroupID:  groupId,
		Topic:    topic,
		MaxBytes: 10e6,
	})
	return r
}
