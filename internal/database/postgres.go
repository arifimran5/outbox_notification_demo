package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Update with your actual connection string
	connStr := os.Getenv("DB_URL")

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening db: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
