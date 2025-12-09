package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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

	// Rate limiter for auth endpoints (5 attempts per minute)
	authLimiter := middleware.NewRateLimiter(5, time.Minute)

	// Admin auth routes (no auth required, but rate limited)
	r.GET("/admin/login", adminCtrl.LoginPage)
	r.POST("/admin/login", middleware.RateLimitMiddleware(authLimiter), adminCtrl.Login)

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

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
