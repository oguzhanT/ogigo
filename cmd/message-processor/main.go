package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"ogigo/internal/processor"
	"ogigo/pkg/queue"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	queueReader := queue.NewKafkaReader("redpanda:29092", "wallet-updates", "wallet-processor")

	connStr := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Printf("Faild to connect to database: %v", err)
	}

	for {
		m, err := queueReader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("error while receiving message: %v", err)
		}

		var event processor.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		processor.ProcessEvent(db, event)

		err2 := queueReader.CommitMessages(context.Background(), m)
		if err2 != nil {
			log.Fatalf("error while committing message: %v", err)
		}

	}
}
