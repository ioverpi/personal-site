package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/controllers"
)

func main() {
	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	r := gin.Default()

	// Static files
	r.Static("/static", "./static")

	// Controllers
	homeCtrl := controllers.NewHomeController()

	// Public routes
	r.GET("/", homeCtrl.Index)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
