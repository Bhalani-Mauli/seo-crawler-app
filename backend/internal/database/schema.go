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
	if err := sm.createCrawlResultsTable(); err != nil {
		return err
	}
	if err := sm.createCrawlLinksTable(); err != nil {
		return err
	}
	if err := sm.createCrawlHeadingsTable(); err != nil {
		return err
	}
	return nil
}

func (sm *SchemaManager) createCrawlResultsTable() error {
	query := `CREATE TABLE IF NOT EXISTS crawl_results (
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
		INDEX idx_created_at (created_at)
	)`
	
	_, err := sm.conn.DB.Exec(query)
	if err != nil {
		log.Printf("Failed to create crawl_results table: %v", err)
		return err
	}
	return nil
}

func (sm *SchemaManager) createCrawlLinksTable() error {
	query := `CREATE TABLE IF NOT EXISTS crawl_links (
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
	)`
	
	_, err := sm.conn.DB.Exec(query)
	if err != nil {
		log.Printf("Failed to create crawl_links table: %v", err)
		return err
	}
	return nil
}

func (sm *SchemaManager) createCrawlHeadingsTable() error {
	query := `CREATE TABLE IF NOT EXISTS crawl_headings (
		id INT AUTO_INCREMENT PRIMARY KEY,
		crawl_result_id INT NOT NULL,
		heading_level VARCHAR(10) NOT NULL,
		heading_text VARCHAR(500),
		heading_order INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (crawl_result_id) REFERENCES crawl_results(id) ON DELETE CASCADE,
		INDEX idx_crawl_result_id (crawl_result_id),
		INDEX idx_heading_level (heading_level)
	)`
	
	_, err := sm.conn.DB.Exec(query)
	if err != nil {
		log.Printf("Failed to create crawl_headings table: %v", err)
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