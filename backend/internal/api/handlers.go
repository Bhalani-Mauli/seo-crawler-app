package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seo-crawler-app/internal/database"
)

// Handler defines the interface for API handlers
type Handler interface {
	HealthCheck(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

// handler implements the Handler interface
type handler struct {
	authHandler  *AuthHandler
	migrationManager *database.MigrationManager
}

// NewHandler creates a new API handler
func NewHandler(authHandler *AuthHandler, migrationManager *database.MigrationManager) Handler {
	return &handler{
		authHandler:  authHandler,
		migrationManager: migrationManager,
	}
}

// HealthCheck handles health check requests
func (h *handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": "2024-01-01T00:00:00Z", 
	})
}

// Register handles user registration
func (h *handler) Register(c *gin.Context) {
	h.authHandler.Register(c)
}

// Login handles user login
func (h *handler) Login(c *gin.Context) {
	h.authHandler.Login(c)
}

// GetProfile returns the current user's profile
func (h *handler) GetProfile(c *gin.Context) {
	h.authHandler.GetProfile(c)
}
