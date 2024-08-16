package storage_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"ogigo/internal/storage"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Load environment variables from the root .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

func setupTestDB() *sql.DB {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return db
}

func TestGetWallet(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	storage := storage.NewStorage(db)

	wallet, err := storage.GetWallet()
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.NotEmpty(t, wallet.WalletID)
}
