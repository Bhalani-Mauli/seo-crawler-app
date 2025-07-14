package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/seo-crawler-app/internal/database"
	"github.com/seo-crawler-app/internal/models"
)

// AuthService handles authentication-related operations
type AuthService struct {
	userRepo *database.UserRepository
	jwtSecret string
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *database.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *AuthService) LoginUser(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return int(userID), nil
	}

	return 0, errors.New("invalid token")
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// generateJWT generates a JWT token for a user
func (s *AuthService) generateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
} 