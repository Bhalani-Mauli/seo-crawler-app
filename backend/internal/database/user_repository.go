package database

import (
	"database/sql"
	"time"

	"github.com/seo-crawler-checker/internal/models"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password, first_name, last_name, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	user.ID = int(id)
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, created_at, updated_at 
			  FROM users WHERE email = ?`
	
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, created_at, updated_at 
			  FROM users WHERE id = ?`
	
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET first_name = ?, last_name = ?, updated_at = ? WHERE id = ?`
	
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UpdatedAt, user.ID)
	return err
}

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
} 