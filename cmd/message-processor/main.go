package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"ogigo/internal/processor"
	"ogigo/internal/queue"
	"os"

	_ "github.com/lib/pq"
)

var kafkaQueue queue.QueueInterface

func main() {

	kafkaQueue = queue.NewQueue("redpanda:29092", "wallet-updates")

	connStr := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Printf("Faild to connect to database: %v", err)
	}

	for {
		m, err := kafkaQueue.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("error while receiving message: %v", err)
		}

		var event processor.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		processor.ProcessEvent(db, event)

		err2 := kafkaQueue.CommitMessages(context.Background(), m)
		if err2 != nil {
			log.Fatalf("error while committing message: %v", err)
		}

	}
}
