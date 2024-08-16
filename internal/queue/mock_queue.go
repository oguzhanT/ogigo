package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MockQueue struct {
	WriteMessageFunc   func(msg kafka.Message) error
	ReadMessageFunc    func(ctx context.Context) (kafka.Message, error)
	CommitMessagesFunc func(ctx context.Context, msgs ...kafka.Message) error
}

func (m *MockQueue) WriteMessage(msg kafka.Message) error {
	return m.WriteMessageFunc(msg)
}

func (m *MockQueue) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return m.ReadMessageFunc(ctx)
}

func (m *MockQueue) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return m.CommitMessagesFunc(ctx, msgs...)
}
