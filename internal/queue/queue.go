package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// QueueInterface abstracts the methods for interacting with Kafka
type QueueInterface interface {
	WriteMessage(msg kafka.Message) error
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

// Queue represents a Kafka queue with a writer and reader
type Queue struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

// NewQueue creates a new Queue instance
func NewQueue(broker, topic string) *Queue {
	return &Queue{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: "wallet-updates",
		}),
	}
}

// WriteMessage writes a message to the Kafka queue
func (q *Queue) WriteMessage(msg kafka.Message) error {
	return q.Writer.WriteMessages(context.Background(), msg)
}

// ReadMessage reads a message from the Kafka queue
func (q *Queue) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return q.Reader.ReadMessage(ctx)
}

// CommitMessages commits the messages to Kafka
func (q *Queue) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return q.Reader.CommitMessages(ctx, msgs...)
}
