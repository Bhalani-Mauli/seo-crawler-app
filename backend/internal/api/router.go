package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seo-crawler-app/internal/config"
	"github.com/seo-crawler-app/internal/middleware"
	"github.com/seo-crawler-app/internal/services"
	"github.com/seo-crawler-app/internal/database"
)	

// Router sets up the HTTP routes
type Router struct {
	handler     Handler
	config      *config.Config
	authService *services.AuthService
	migrationManager *database.MigrationManager
}

// NewRouter creates a new router instance
func NewRouter(handler Handler, config *config.Config, authService *services.AuthService, migrationManager *database.MigrationManager) *Router {
	return &Router{
		handler:     handler,
		config:      config,
		authService: authService,
		migrationManager: migrationManager,
	}
}

// SetupRoutes configures all the routes
func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}))


	// Health check endpoint (no auth required)
	router.GET("/ping", r.handler.HealthCheck)
	
	// Public auth routes (no auth required)
	router.POST("/api/auth/register", r.handler.Register)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(r.authService))

	router.OPTIONS("/*path", func(c *gin.Context) {
	    c.Status(204)
	})

	return router
}