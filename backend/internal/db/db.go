package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initialises the global database connection pool
func Init(databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return nil
}

// Close cleanly shuts down the database pool
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func RunMigrations() error {
    path := "internal/db/migrations/init.sql"
    
    c, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("could not read migration file: %w", err)
    }

    _, err = DB.Exec(string(c))
    if err != nil {
        return fmt.Errorf("could not execute migration: %w", err)
    }

    return nil
}