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

	// Execute the query
	row := s.DB.QueryRow(query, walletId, currency)

	// Initialize a Wallet struct
	var wallet Wallet

	// Scan the result into the wallet struct
	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result found
			log.Printf("No wallet found with ID: %v and Currency: %v", walletId, currency)
			return Wallet{}, nil // or return an empty wallet with initial values
		}
		return Wallet{}, fmt.Errorf("failed to scan wallet: %w", err)
	}

	// Return the wallet found
	return wallet, nil
}

// UpdateBalance updates the wallet balance in the database
func (s *Storage) UpdateBalance(walletID string, amount float64, currency string) error {
	query := `
        INSERT INTO wallets (wallet_id, balance, currency)
        VALUES ($1, $2, $3)
        ON CONFLICT (wallet_id, currency)
        DO UPDATE SET balance = wallets.balance + $2
    `
	_, err := s.DB.Exec(query, walletID, amount, currency)
	if err != nil {
		log.Printf(query, walletID, amount, currency)
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}
