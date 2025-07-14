// main.go
package main

import (
	"log"

	"github.com/seo-crawler-app/internal/api"
	"github.com/seo-crawler-app/internal/config"
	"github.com/seo-crawler-app/internal/database"
	"github.com/seo-crawler-app/internal/services"
	"github.com/seo-crawler-app/pkg/crawler"
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

	crawlRepo := database.NewCrawlRepository(dbConn)
	userRepo := database.NewUserRepository(dbConn.DB)

	crawlerService := crawler.NewCrawler()

	crawlService := services.NewCrawlService(crawlRepo, crawlerService)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)

	authHandler := api.NewAuthHandler(authService)
	handler := api.NewHandler(crawlService, authHandler, migrationManager)

	router := api.NewRouter(handler, cfg, authService, migrationManager)
	app := router.SetupRoutes()

	log.Printf("Server starting on :%s", cfg.Server.Port)
	if err := app.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
