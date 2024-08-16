package storage

import (
	"database/sql"
	"fmt"
	"log"
)

type Storage struct {
	DB *sql.DB
}
type Wallet struct {
	WalletID string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func (s *Storage) GetWallet() (Wallet, error) {
	// Test values
	walletId := "01HPMV01XPAXCG242W7SZWD0S5"
	currency := "TRY"

	query := `
		SELECT wallet_id, balance, currency
		FROM wallets 
		WHERE wallet_id = $1 AND currency = $2
	`
	row := s.DB.QueryRow(query, walletId, currency)

	var wallet Wallet

	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No wallet found with ID: %v and Currency: %v", walletId, currency)
			return Wallet{}, nil // or return an empty wallet with initial values
		}
		return Wallet{}, fmt.Errorf("failed to scan wallet: %w", err)
	}

	return wallet, nil
}

// UpdateBalance updates the wallet balance in the database
func (s *Storage) UpdateWalletBalance(walletID string, amount float64, currency string) error {
	// Start a new transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	selectQuery := `
        SELECT balance
        FROM wallets
        WHERE wallet_id = $1 AND currency = $2
        FOR UPDATE
    `
	// Execute the select query to lock the row
	row := tx.QueryRow(selectQuery, walletID, currency)
	var currentBalance float64
	err = row.Scan(&currentBalance)
	if err != nil && err != sql.ErrNoRows {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("Failed to rollback transaction: %v", err2)
		}

		log.Printf("Failed to select row: %v", err)
		return fmt.Errorf("failed to select row: %w", err)
	}

	// Define the update query
	updateQuery := `
        INSERT INTO wallets (wallet_id, balance, currency)
        VALUES ($1, $2, $3)
        ON CONFLICT (wallet_id, currency)
        DO UPDATE SET balance = wallets.balance + $2
    `
	_, err = tx.Exec(updateQuery, walletID, amount, currency)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("Failed to rollback transaction: %v", err2)
		}

		log.Printf("Failed to execute update query: %s, params: %v, %v, %v", updateQuery, walletID, amount, currency)
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
