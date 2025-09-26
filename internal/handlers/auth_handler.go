package handlers

import (
	"log"
	"net/http"
	"wallet/internal/models"
	"wallet/pkg/auth"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: input.Password,
	}

	if err := user.Create(); err != nil {
		log.Printf("error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create a wallet for the new user
	wallet := models.Wallet{}
	if err := wallet.Create(user.ID); err != nil {
		log.Printf("error creating wallet for user %s: %v", user.ID, err)
		// Handle wallet creation failure - maybe delete the user or log for manual intervention
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet for user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
