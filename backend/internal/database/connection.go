package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Connection represents a database connection
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new database connection
func NewConnection(dsn string) (*Connection, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to MySQL!")
	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.DB.Close()
} 