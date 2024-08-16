package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ogigo/internal/queue"
	"ogigo/internal/storage"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var kafkaQueue queue.QueueInterface

type Event struct {
	App  string `json:"app"`
	Type string `json:"type"`
	Time string `json:"time"`
	Meta struct {
		User string `json:"user"`
	} `json:"meta"`
	Wallet     string `json:"wallet"`
	Attributes struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"attributes"`
}

type EventRequest struct {
	Events []Event `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Wallet struct {
	ID       string    `json:"id"`
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type WalletResponse struct {
	Wallets []Wallet `json:"wallets"`
}

// @title           Ogigo API
// @version         1.0
// @description     This is a test scenario server.

// @host      localhost:80
// @BasePath  /
func main() {
	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	// Initialize Kafka queue
	kafkaQueue = queue.NewQueue("redpanda:29092", "wallet-updates")

	r := gin.Default()

	// Middleware to capture Gin errors
	r.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			sentry.CaptureException(c.Errors.Last())
		}
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	r.POST("/", handleUpdate)
	r.GET("/", handleState)

	if err := r.Run(":80"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// handleUpdate processes incoming events and places them into the queue.
// @Summary      Process incoming wallet update events
// @Description  This endpoint processes wallet update events and places them into a Kafka queue.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        events  body      EventRequest  true  "Events to be processed"
// @Success      200     {string}  string        "OK"
// @Failure      400     {object}  ErrorResponse "Bad Request"
// @Failure      500     {object}  ErrorResponse "Internal Server Error"
// @Router       / [post]
func handleUpdate(c *gin.Context) {
	var req EventRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, event := range req.Events {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			errorResponse := ErrorResponse{Error: "Failed to marshal event"}
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}

		err = kafkaQueue.WriteMessage(kafka.Message{
			Value: eventJSON,
		})
		if err != nil {
			errorResponse := ErrorResponse{Error: "Failed to write message to queue"}
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}
	}

	c.JSON(http.StatusOK, "OK")
}

// handleState retrieves the current wallet state.
// @Summary      Retrieve the current wallet state
// @Description  This endpoint retrieves the current wallet state.
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Success      200     {object}  WalletResponse "OK"
// @Failure      500     {object}  ErrorResponse   "Internal Server Error"
// @Router       / [get]
func handleState(c *gin.Context) {
	connStr := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable"

	log.Printf("Connecting to the database with connection string: %v", connStr)
	// Attempt to open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	// Initialize the storage
	storage := storage.NewStorage(db)

	// Get the wallet data
	wallet, err := storage.GetWallet()
	if err != nil {
		errorResponse := ErrorResponse{Error: "Failed to get wallet data"}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Check if wallet is empty (in case no record was found)
	if wallet.WalletID == "" {
		errorResponse := ErrorResponse{Error: "No wallet found with the specified criteria"}
		c.JSON(http.StatusOK, errorResponse)
		return
	}

	// Prepare the response
	walletResponse := WalletResponse{
		Wallets: []Wallet{
			{
				ID: wallet.WalletID,
				Balances: []Balance{
					{
						Currency: wallet.Currency,
						Amount:   wallet.Balance,
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, walletResponse)
}
