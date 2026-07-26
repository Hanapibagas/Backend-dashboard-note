package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"auth/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes database connection
func InitDB(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	DB = db

	log.Println("Database connected successfully")
	return nil
}

// CloseDB closes database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns database instance
func GetDB() *sql.DB {
	return DB
}
