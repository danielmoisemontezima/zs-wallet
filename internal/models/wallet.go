package models

import (
	"context"
	"time"
	"wallet/pkg/database"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w *Wallet) Create(userID uuid.UUID) error {
	w.ID = uuid.New()
	w.UserID = userID
	w.Balance = 0.00
	w.Currency = "USD"
	query := `INSERT INTO wallets (id, user_id, balance, currency) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`
	err := database.DB.QueryRow(context.Background(), query, w.ID, w.UserID, w.Balance, w.Currency).Scan(&w.CreatedAt, &w.UpdatedAt)
	return err
}

func GetWalletByUserID(userID uuid.UUID) (*Wallet, error) {
	var wallet Wallet
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM wallets WHERE user_id = $1`
	err := database.DB.QueryRow(context.Background(), query, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
