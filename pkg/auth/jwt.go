package auth

import (
	"errors"
	"time"
	"wallet/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func Initialize(cfg *config.Config) {
	jwtKey = []byte(cfg.JWTSecret)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("could not parse token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
