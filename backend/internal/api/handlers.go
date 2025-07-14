package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seo-crawler-app/internal/models"
	"github.com/seo-crawler-app/internal/services"
	"github.com/seo-crawler-app/internal/database"
)

// Handler defines the interface for API handlers
type Handler interface {
	SubmitCrawl(c *gin.Context)
	GetResults(c *gin.Context)
	GetResultByID(c *gin.Context)
	GetLinksByID(c *gin.Context)
	GetHeadingsByID(c *gin.Context)
	BulkRerun(c *gin.Context)
	BulkDelete(c *gin.Context)
	StopCrawl(c *gin.Context)
	HealthCheck(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	GetMigrationStatus(c *gin.Context)
}

// handler implements the Handler interface
type handler struct {
	crawlService services.CrawlService
	authHandler  *AuthHandler
	migrationManager *database.MigrationManager
}

// NewHandler creates a new API handler
func NewHandler(crawlService services.CrawlService, authHandler *AuthHandler, migrationManager *database.MigrationManager) Handler {
	return &handler{
		crawlService: crawlService,
		authHandler:  authHandler,
		migrationManager: migrationManager,
	}
}

// SubmitCrawl handles crawl submission requests
func (h *handler) SubmitCrawl(c *gin.Context) {
	var req models.CrawlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}


	if err := h.crawlService.SubmitCrawl(userID.(int), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL submitted for crawling"})
}

// GetResults handles paginated results retrieval
func (h *handler) GetResults(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	response, err := h.crawlService.GetCrawlResults(userID.(int), page, pageSize, status, search, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetResultByID handles single result retrieval by ID
func (h *handler) GetResultByID(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")

	result, err := h.crawlService.GetCrawlResultByID(userID.(int), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLinksByID handles links retrieval for a specific crawl result
func (h *handler) GetLinksByID(c *gin.Context) {
	id := c.Param("id")

	links, err := h.crawlService.GetLinksByCrawlID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}

// GetHeadingsByID handles headings retrieval for a specific crawl result
func (h *handler) GetHeadingsByID(c *gin.Context) {
	id := c.Param("id")

	headings, err := h.crawlService.GetHeadingsByCrawlID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch headings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"headings": headings})
}

// BulkRerun handles bulk re-crawl requests
func (h *handler) BulkRerun(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.BulkActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.crawlService.BulkRerun(userID.(int), req.URLs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Re-crawling URLs"})
}

// BulkDelete handles bulk deletion requests
func (h *handler) BulkDelete(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.BulkActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.crawlService.BulkDelete(userID.(int), req.URLs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URLs deleted successfully"})
}

// StopCrawl handles crawl stopping requests
func (h *handler) StopCrawl(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")

	if err := h.crawlService.StopCrawl(userID.(int), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop crawling"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crawling stopped"})
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

// GetMigrationStatus returns the status of all migrations
func (h *handler) GetMigrationStatus(c *gin.Context) {
	status, err := h.migrationManager.GetMigrationStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get migration status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"migrations": status,
		"total":      len(status),
	})
} 