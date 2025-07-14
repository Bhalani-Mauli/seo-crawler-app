// main.go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/seo-crawler-checker/internal/config"
	"github.com/seo-crawler-checker/internal/database"
)

func main() {
	cfg := config.Load()

	dbConn, err := database.NewConnection(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	migrationManager := database.NewMigrationManager(dbConn)
	if err := migrationManager.Initialize(); err != nil {
		log.Fatal("Failed to initialize database migrations:", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080") // Start server on port 8080
}
