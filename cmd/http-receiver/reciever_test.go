package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafkaWriter is a mock implementation of queue.Writer for testing.
type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestHandleUpdate(t *testing.T) {
	mockKafkaWriter := new(MockKafkaWriter)
	queueWriter = mockKafkaWriter

	r := gin.Default()
	r.POST("/", handleUpdate)

	tests := []struct {
		name           string
		requestBody    EventRequest
		mockWriteError error
		expectedCode   int
		expectedBody   string
	}{
		{
			name: "successful write",
			requestBody: EventRequest{
				Events: []Event{
					{
						App:  "test-app",
						Type: "update",
						Time: "2024-08-16T12:00:00Z",
						Meta: struct {
							User string `json:"user"`
						}{
							User: "test-user",
						},
						Wallet: "test-wallet",
						Attributes: struct {
							Amount   string `json:"amount"`
							Currency string `json:"currency"`
						}{
							Amount:   "100",
							Currency: "USD",
						},
					},
				},
			},
			mockWriteError: nil,
			expectedCode:   http.StatusOK,
			expectedBody:   `"OK"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(tt.mockWriteError)

			body, _ := json.Marshal(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}
