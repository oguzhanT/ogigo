package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Writer is an interface for Kafka Writer
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

// NewKafkaWriter creates a new Kafka Writer
func NewKafkaWriter(broker, topic string) Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewKafkaReader(broker, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
}
