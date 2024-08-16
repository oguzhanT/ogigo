package processor

import (
	"database/sql"
	"fmt"
	"ogigo/internal/storage"
	"strconv"
)

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

func ProcessEvent(db *sql.DB, event Event) {
	fmt.Printf("Processing event: %+v\n", event)

	amount, err := strconv.ParseFloat(event.Attributes.Amount, 64)
	if err != nil {
		fmt.Printf("Failed to parse amount: %v", err)
		return
	}

	currency := event.Attributes.Currency

	var balanceChange float64
	if event.Type == "BALANCE_INCREASE" {
		balanceChange = amount
	} else if event.Type == "BALANCE_DECREASE" {
		balanceChange = -amount
	} else {
		fmt.Printf("Unknown event type: %s", event.Type)
		return
	}

	storage := storage.NewStorage(db)
	err = storage.UpdateWalletBalance(event.Wallet, balanceChange, currency)
	if err != nil {
		fmt.Printf("Failed to update balance: %v", err)
	}
}
