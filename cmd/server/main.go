package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/controllers"
	"github.com/ioverpi/personal-site/internal/services"
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

	// Services
	blogService := services.NewBlogService(application)
	projectsService := services.NewProjectsService(application)
	quotesService := services.NewQuotesService(application)

	// Controllers
	homeCtrl := controllers.NewHomeController()
	blogCtrl := controllers.NewBlogController(blogService)
	projectsCtrl := controllers.NewProjectsController(projectsService)
	quotesCtrl := controllers.NewQuotesController(quotesService)

	// Public routes
	r.GET("/", homeCtrl.Index)
	r.GET("/blog", blogCtrl.List)
	r.GET("/blog/:slug", blogCtrl.Show)
	r.GET("/projects", projectsCtrl.List)
	r.GET("/quotes", quotesCtrl.List)
	r.GET("/quotes/random", quotesCtrl.Random)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
