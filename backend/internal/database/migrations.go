package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Migration represents a database migration
type Migration struct {
	ID          int
	Name        string
	SQL         string
	Description string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	conn *Connection
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(conn *Connection) *MigrationManager {
	return &MigrationManager{conn: conn}
}

// Initialize creates the migrations table and runs pending migrations
func (mm *MigrationManager) Initialize() error {
	// Create migrations table if it doesn't exist
	if err := mm.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migrations
	migrations := mm.getMigrations()

	// Run pending migrations
	if err := mm.runPendingMigrations(migrations); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (mm *MigrationManager) createMigrationsTable() error {
	query := `CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_name (name)
	)`
	
	_, err := mm.conn.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	return nil
}

// getMigrations returns all available migrations in order
func (mm *MigrationManager) getMigrations() []Migration {
	return []Migration{
		{
			ID:          1,
			Name:        "001_create_users_table",
			Description: "Create users table for authentication",
			SQL: `CREATE TABLE IF NOT EXISTS users (
				id INT AUTO_INCREMENT PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				password VARCHAR(255) NOT NULL,
				first_name VARCHAR(100) NOT NULL,
				last_name VARCHAR(100) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				INDEX idx_email (email)
			)`,
		},
		{
			ID:          2,
			Name:        "002_create_crawl_results_table",
			Description: "Create crawl_results table with user_id",
			SQL: `CREATE TABLE IF NOT EXISTS crawl_results (
				id INT AUTO_INCREMENT PRIMARY KEY,
				user_id INT NOT NULL,
				url VARCHAR(500) NOT NULL,
				html_version VARCHAR(50),
				title VARCHAR(500),
				headings JSON,
				internal_links INT DEFAULT 0,
				external_links INT DEFAULT 0,
				inaccessible_links INT DEFAULT 0,
				has_login_form BOOLEAN DEFAULT FALSE,
				status VARCHAR(50) DEFAULT 'pending',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				INDEX idx_user_id (user_id),
				INDEX idx_status (status),
				INDEX idx_created_at (created_at),
				UNIQUE KEY unique_user_url (user_id, url)
			)`,
		},
		{
			ID:          3,
			Name:        "003_create_crawl_links_table",
			Description: "Create crawl_links table for storing link details",
			SQL: `CREATE TABLE IF NOT EXISTS crawl_links (
				id INT AUTO_INCREMENT PRIMARY KEY,
				crawl_result_id INT NOT NULL,
				link_url VARCHAR(1000) NOT NULL,
				link_text VARCHAR(500),
				link_type ENUM('internal', 'external') NOT NULL,
				status_code INT,
				is_accessible BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (crawl_result_id) REFERENCES crawl_results(id) ON DELETE CASCADE,
				INDEX idx_crawl_result_id (crawl_result_id),
				INDEX idx_link_type (link_type),
				INDEX idx_is_accessible (is_accessible)
			)`,
		},
		{
			ID:          4,
			Name:        "004_create_crawl_headings_table",
			Description: "Create crawl_headings table for storing heading details",
			SQL: `CREATE TABLE IF NOT EXISTS crawl_headings (
				id INT AUTO_INCREMENT PRIMARY KEY,
				crawl_result_id INT NOT NULL,
				heading_level VARCHAR(10) NOT NULL,
				heading_text VARCHAR(500),
				heading_order INT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (crawl_result_id) REFERENCES crawl_results(id) ON DELETE CASCADE,
				INDEX idx_crawl_result_id (crawl_result_id),
				INDEX idx_heading_level (heading_level)
			)`,
		},
	}
}

// runPendingMigrations runs all migrations that haven't been applied yet
func (mm *MigrationManager) runPendingMigrations(migrations []Migration) error {
	for _, migration := range migrations {
		applied, err := mm.isMigrationApplied(migration.Name)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if applied {
			log.Printf("Migration %s already applied, skipping", migration.Name)
			continue
		}

		log.Printf("Applying migration: %s - %s", migration.Name, migration.Description)
		
		tx, err := mm.conn.DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(migration.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Name, err)
		}

		if err := mm.recordMigrationApplied(tx, migration); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Name, err)
		}

		log.Printf("Successfully applied migration: %s", migration.Name)
	}

	return nil
}

// isMigrationApplied checks if a migration has already been applied
func (mm *MigrationManager) isMigrationApplied(name string) (bool, error) {
	var count int
	err := mm.conn.DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigrationApplied records that a migration has been applied
func (mm *MigrationManager) recordMigrationApplied(tx *sql.Tx, migration Migration) error {
	query := `INSERT INTO migrations (name, description, applied_at) VALUES (?, ?, ?)`
	_, err := tx.Exec(query, migration.Name, migration.Description, time.Now())
	return err
}

// GetMigrationStatus returns the status of all migrations
func (mm *MigrationManager) GetMigrationStatus() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			m.name,
			m.description,
			m.applied_at,
			CASE WHEN m.name IS NOT NULL THEN 'applied' ELSE 'pending' END as status
		FROM migrations m
		ORDER BY m.applied_at
	`
	
	rows, err := mm.conn.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statusList []map[string]interface{}
	for rows.Next() {
		var name, description, status string
		var appliedAt sql.NullTime
		
		if err := rows.Scan(&name, &description, &appliedAt, &status); err != nil {
			return nil, err
		}

		statusMap := map[string]interface{}{
			"name":        name,
			"description": description,
			"status":      status,
		}
		
		if appliedAt.Valid {
			statusMap["applied_at"] = appliedAt.Time
		}

		statusList = append(statusList, statusMap)
	}

	return statusList, nil
} 