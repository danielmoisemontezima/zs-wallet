package main

import (
	"fmt"
	"log"
	"wallet/config"
	"wallet/internal/handlers"
	"wallet/internal/middleware"
	"wallet/pkg/auth"
	"wallet/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.Connect(cfg)
	database.RunMigrations(database.DB)
	auth.Initialize(cfg)

	r := gin.Default()

	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", handlers.Register)
			authRoutes.POST("/login", handlers.Login)
		}

		walletRoutes := api.Group("/wallet")
		walletRoutes.Use(middleware.AuthMiddleware())
		{
			walletRoutes.GET("/balance", handlers.GetBalance)
			walletRoutes.POST("/deposit", handlers.Deposit)
			walletRoutes.POST("/withdraw", handlers.Withdraw)
			walletRoutes.GET("/transactions", handlers.GetTransactions)
		}
	}

	log.Printf("Server starting on port %s", cfg.APIPort)
	r.Run(fmt.Sprintf(":%s", cfg.APIPort))
}
