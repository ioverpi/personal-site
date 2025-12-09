package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/controllers"
	"github.com/ioverpi/personal-site/internal/database"
	"github.com/ioverpi/personal-site/internal/middleware"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/migrations"
)


func main() {
	cfg := config.Load()

	// Set migrations for database package
	database.MigrationsFS = migrations.FS

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	r := gin.Default()

	// Security headers
	r.Use(middleware.SecurityHeaders())

	// Static files
	r.Static("/static", "./static")

	// Services
	blogService := services.NewBlogService(application)
	projectsService := services.NewProjectsService(application)
	quotesService := services.NewQuotesService(application)
	adminService := services.NewAdminService(application)
	authService := services.NewAuthService(application)
	userService := services.NewUserService(application)

	// Controllers
	homeCtrl := controllers.NewHomeController()
	blogCtrl := controllers.NewBlogController(blogService)
	projectsCtrl := controllers.NewProjectsController(projectsService)
	quotesCtrl := controllers.NewQuotesController(quotesService)
	adminCtrl := controllers.NewAdminController(
		adminService,
		blogService,
		projectsService,
		quotesService,
		authService,
		userService,
		cfg,
	)

	// Public routes
	r.GET("/", homeCtrl.Index)
	r.GET("/blog", blogCtrl.List)
	r.GET("/blog/:slug", blogCtrl.Show)
	r.GET("/projects", projectsCtrl.List)
	r.GET("/quotes", quotesCtrl.List)
	r.GET("/quotes/random", quotesCtrl.Random)

	// Registration (public, via invite token)
	r.GET("/register", adminCtrl.RegisterPage)
	r.POST("/register", adminCtrl.Register)

	// Admin auth routes (no auth required)
	r.GET("/admin/login", adminCtrl.LoginPage)
	r.POST("/admin/login", adminCtrl.Login)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService, cfg.SecureCookies))
	{
		admin.GET("/logout", adminCtrl.Logout)
		admin.GET("/", adminCtrl.Dashboard)

		// Users (admin only)
		admin.GET("/users", adminCtrl.UsersList)
		admin.GET("/invites/new", adminCtrl.NewInvite)
		admin.POST("/invites", adminCtrl.CreateInvite)
		admin.POST("/invites/:id/delete", adminCtrl.DeleteInvite)

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
