package models

import (
	"context"
	"time"
	"wallet/pkg/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Create() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.ID = uuid.New()
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING created_at, updated_at`
	err = database.DB.QueryRow(context.Background(), query, u.ID, u.Email, string(hashedPassword)).Scan(&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := database.DB.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
