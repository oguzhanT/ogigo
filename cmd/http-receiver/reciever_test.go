package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ogigo/internal/queue"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

// Define a mock queue for testing
type MockQueue struct {
	queue.QueueInterface
	WriteMessageFunc   func(msg kafka.Message) error
	ReadMessageFunc    func(ctx context.Context) (kafka.Message, error)
	CommitMessagesFunc func(ctx context.Context, msgs ...kafka.Message) error
}

func (mq *MockQueue) WriteMessage(msg kafka.Message) error {
	return mq.WriteMessageFunc(msg)
}

func (mq *MockQueue) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return mq.ReadMessageFunc(ctx)
}

func (mq *MockQueue) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return mq.CommitMessagesFunc(ctx, msgs...)
}

func TestHandleUpdate_successful_write(t *testing.T) {
	mockQueue := &MockQueue{
		WriteMessageFunc: func(msg kafka.Message) error {
			return nil
		},
	}
	kafkaQueue = mockQueue // Assign the mockQueue to kafkaQueue

	r := gin.Default()
	r.POST("/", handleUpdate)

	event := Event{
		App:  "test-app",
		Type: "update",
		Time: "2024-08-16T12:00:00Z",
		Meta: struct {
			User string "json:\"user\""
		}{
			User: "test-user",
		},
		Wallet: "test-wallet",
		Attributes: struct {
			Amount   string "json:\"amount\""
			Currency string "json:\"currency\""
		}{
			Amount:   "100",
			Currency: "USD",
		},
	}
	reqBody, _ := json.Marshal(EventRequest{Events: []Event{event}})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `"OK"`, w.Body.String())
}

func TestHandleUpdate_failed_to_write(t *testing.T) {
	mockQueue := &MockQueue{
		WriteMessageFunc: func(msg kafka.Message) error {
			return fmt.Errorf("write error")
		},
	}
	kafkaQueue = mockQueue // Assign the mockQueue to kafkaQueue

	r := gin.Default()
	r.POST("/", handleUpdate)

	event := Event{
		App:  "test-app",
		Type: "update",
		Time: "2024-08-16T12:00:00Z",
		Meta: struct {
			User string "json:\"user\""
		}{
			User: "test-user",
		},
		Wallet: "test-wallet",
		Attributes: struct {
			Amount   string "json:\"amount\""
			Currency string "json:\"currency\""
		}{
			Amount:   "100",
			Currency: "USD",
		},
	}
	reqBody, _ := json.Marshal(EventRequest{Events: []Event{event}})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Failed to write message to queue"}`, w.Body.String())
}
