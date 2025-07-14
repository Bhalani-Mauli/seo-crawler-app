package database

import (
	"log"
)

// SchemaManager handles database schema operations
type SchemaManager struct {
	conn *Connection
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(conn *Connection) *SchemaManager {
	return &SchemaManager{conn: conn}
}

// Initialize creates all necessary tables
func (sm *SchemaManager) Initialize() error {
	if err := sm.createUsersTable(); err != nil {
		return err
	}
	return nil
}

func (sm *SchemaManager) createUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_email (email)
	)`
	
	_, err := sm.conn.DB.Exec(query)
	if err != nil {
		log.Printf("Failed to create users table: %v", err)
		return err
	}
	return nil
}