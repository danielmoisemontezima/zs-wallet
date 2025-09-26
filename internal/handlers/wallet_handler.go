package handlers

import (
	"context"
	"errors"
	"net/http"
	"wallet/internal/models"
	"wallet/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetBalance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	wallet, err := models.GetWalletByUserID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

type DepositInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func Deposit(c *gin.Context) {
	var input DepositInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	userUUID, _ := uuid.Parse(userID.(string))

	err := performTransaction(userUUID, "deposit", input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

type WithdrawInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func Withdraw(c *gin.Context) {
	var input WithdrawInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	userUUID, _ := uuid.Parse(userID.(string))

	err := performTransaction(userUUID, "withdraw", input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
}

func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID, _ := uuid.Parse(userID.(string))

	wallet, err := models.GetWalletByUserID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	transactions, err := models.GetTransactionsByWalletID(wallet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func performTransaction(userID uuid.UUID, txType string, amount float64) error {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer tx.Rollback(context.Background())

	// Get wallet with row-level lock
	wallet, err := getWalletForUpdate(tx, userID)
	if err != nil {
		return errors.New("wallet not found")
	}

	balanceBefore := wallet.Balance
	var balanceAfter float64
	description := ""

	switch txType {
	case "deposit":
		balanceAfter = balanceBefore + amount
		description = "User deposit"
	case "withdraw":
		if balanceBefore < amount {
			return errors.New("insufficient funds")
		}
		balanceAfter = balanceBefore - amount
		description = "User withdrawal"
	default:
		return errors.New("invalid transaction type")
	}

	// Update wallet balance
	err = updateWalletBalance(tx, wallet.ID, balanceAfter)
	if err != nil {
		return errors.New("failed to update wallet balance")
	}

	// Create transaction record
	err = models.CreateTransaction(tx, wallet.ID, txType, amount, balanceBefore, balanceAfter, description)
	if err != nil {
		return errors.New("failed to create transaction record")
	}

	return tx.Commit(context.Background())
}

func getWalletForUpdate(tx pgx.Tx, userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM wallets WHERE user_id = $1 FOR UPDATE`
	err := tx.QueryRow(context.Background(), query, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func updateWalletBalance(tx pgx.Tx, walletID uuid.UUID, newBalance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.Exec(context.Background(), query, newBalance, walletID)
	return err
}
