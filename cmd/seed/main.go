package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/database"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/migrations"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: seed <password>")
		fmt.Println("Sets the password for the initial admin user (kellon08@gmail.com)")
		os.Exit(1)
	}

	password := os.Args[1]
	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters")
	}

	cfg := config.Load()

	// Set migrations for database package
	database.MigrationsFS = migrations.FS

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	authService := services.NewAuthService(application)
	userService := services.NewUserService(application)

	// Find the admin user
	user, err := userService.GetByEmail("kellon08@gmail.com")
	if err != nil {
		log.Fatalf("Admin user not found. Run migrations first: %v", err)
	}

	// Check if login already exists
	existing, _ := authService.GetLoginByEmail(user.Email)
	if existing != nil {
		log.Println("Login already exists for admin user. Updating password...")
		if err := authService.UpdatePassword(user.ID, password); err != nil {
			log.Fatalf("Failed to update password: %v", err)
		}
		log.Println("Password updated successfully!")
		return
	}

	// Create login
	_, err = authService.CreatePasswordLogin(user.ID, user.Email, password)
	if err != nil {
		log.Fatalf("Failed to create login: %v", err)
	}

	log.Printf("Admin user password set successfully for %s", user.Email)
}
