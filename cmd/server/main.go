package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/controllers"
	"github.com/ioverpi/personal-site/internal/middleware"
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
	adminService := services.NewAdminService(application)

	// Controllers
	homeCtrl := controllers.NewHomeController()
	blogCtrl := controllers.NewBlogController(blogService)
	projectsCtrl := controllers.NewProjectsController(projectsService)
	quotesCtrl := controllers.NewQuotesController(quotesService)
	adminCtrl := controllers.NewAdminController(adminService, blogService, projectsService, quotesService, cfg.AdminPassword)

	// Public routes
	r.GET("/", homeCtrl.Index)
	r.GET("/blog", blogCtrl.List)
	r.GET("/blog/:slug", blogCtrl.Show)
	r.GET("/projects", projectsCtrl.List)
	r.GET("/quotes", quotesCtrl.List)
	r.GET("/quotes/random", quotesCtrl.Random)

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth(cfg.AdminPassword))
	{
		admin.GET("/login", adminCtrl.LoginPage)
		admin.POST("/login", adminCtrl.Login)
		admin.GET("/logout", adminCtrl.Logout)
		admin.GET("/", adminCtrl.Dashboard)

		// Posts
		admin.GET("/posts/new", adminCtrl.NewPost)
		admin.POST("/posts", adminCtrl.CreatePost)
		admin.GET("/posts/:id/edit", adminCtrl.EditPost)
		admin.POST("/posts/:id", adminCtrl.UpdatePost)
		admin.POST("/posts/:id/delete", adminCtrl.DeletePost)

		// Projects
		admin.GET("/projects/new", adminCtrl.NewProject)
		admin.POST("/projects", adminCtrl.CreateProject)
		admin.GET("/projects/:id/edit", adminCtrl.EditProject)
		admin.POST("/projects/:id", adminCtrl.UpdateProject)
		admin.POST("/projects/:id/delete", adminCtrl.DeleteProject)

		// Quotes
		admin.GET("/quotes/new", adminCtrl.NewQuote)
		admin.POST("/quotes", adminCtrl.CreateQuote)
		admin.GET("/quotes/:id/edit", adminCtrl.EditQuote)
		admin.POST("/quotes/:id", adminCtrl.UpdateQuote)
		admin.POST("/quotes/:id/delete", adminCtrl.DeleteQuote)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
