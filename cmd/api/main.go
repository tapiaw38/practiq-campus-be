package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/config"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/database"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/storage"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases"
)

func main() {
	loadConfig()

	cfg := config.GetConfigService()
	gin.SetMode(cfg.ServerConfig.GinMode)

	db, err := database.GetSQLClient()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	ds := datasources.CreateDatasources(db)
	reposFactory := repositories.NewFactory(ds)
	repos := reposFactory()
	integ := integrations.CreateIntegrations(cfg.ServerConfig.AuthAPIURL, cfg.ServerConfig.PractiqAPIURL)
	store := storage.NewS3Storage(cfg.S3Config)
	factory := appcontext.NewFactory(repos, integ, store)
	uc := usecases.NewUsecases(factory)

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "practiq-campus-be"})
	})

	web.RegisterRoutes(app, uc, repos.Profile)

	port := cfg.ServerConfig.Port
	log.Printf("practiq-campus-be running on port %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
