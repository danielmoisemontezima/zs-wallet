package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"wallet/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(cfg *config.Config) {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Database connection established")
}

func RunMigrations(db *pgxpool.Pool) {
	log.Println("Running database migrations...")
	content, err := os.ReadFile("migrations/0001_init.sql")
	if err != nil {
		log.Fatalf("Could not read migration file: %v", err)
	}

	// Simple parser for up migrations
	upSQL := ""
	inUp := false
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "-- +migrate Up") {
			inUp = true
			continue
		}
		if strings.Contains(line, "-- +migrate Down") {
			inUp = false
			continue
		}
		if inUp {
			upSQL += line + "\n"
		}
	}

	_, err = db.Exec(context.Background(), upSQL)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully.")
}
