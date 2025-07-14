package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seo-crawler-checker/internal/config"
	"github.com/seo-crawler-checker/internal/database"
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
	
	// Migration status endpoint (no auth required)
	router.GET("/api/migrations/status", r.handler.GetMigrationStatus)

	// Public auth routes (no auth required)
	router.POST("/api/auth/register", r.handler.Register)
	router.POST("/api/auth/login", r.handler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(r.authService))

	// Bulk action routes
	protected.POST("/bulk/rerun", r.handler.BulkRerun)
	protected.DELETE("/bulk/delete", r.handler.BulkDelete)

	// Control routes
	protected.POST("/stop/:id", r.handler.StopCrawl)

	// User profile route
	protected.GET("/profile", r.handler.GetProfile)

	router.OPTIONS("/*path", func(c *gin.Context) {
	    c.Status(204)
	})

	return router
}