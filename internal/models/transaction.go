package models

import (
	"context"
	"time"
	"wallet/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	ID            uuid.UUID `json:"id"`
	WalletID      uuid.UUID `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

func CreateTransaction(tx pgx.Tx, walletID uuid.UUID, transType string, amount, balanceBefore, balanceAfter float64, description string) error {
	query := `INSERT INTO transactions (id, wallet_id, type, amount, balance_before, balance_after, description)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(context.Background(), query, uuid.New(), walletID, transType, amount, balanceBefore, balanceAfter, description)
	return err
}

func GetTransactionsByWalletID(walletID uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	query := `SELECT id, wallet_id, type, amount, balance_before, balance_after, description, created_at 
              FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(context.Background(), query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.WalletID, &t.Type, &t.Amount, &t.BalanceBefore, &t.BalanceAfter, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
