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
	if len(os.Args) < 4 {
		fmt.Println("Usage: seed <email> <name> <password>")
		fmt.Println("Creates the initial admin user")
		os.Exit(1)
	}

	email := os.Args[1]
	name := os.Args[2]
	password := os.Args[3]

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

	// Check if user already exists
	existing, _ := userService.GetByEmail(email)
	if existing != nil {
		log.Println("User already exists. Updating password...")
		if err := authService.UpdatePassword(existing.ID, password); err != nil {
			log.Fatalf("Failed to update password: %v", err)
		}
		log.Println("Password updated successfully!")
		return
	}

	// Create user
	user, err := userService.CreateUser(services.CreateUserInput{
		Email: email,
		Name:  name,
		Role:  "admin",
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create login
	_, err = authService.CreatePasswordLogin(user.ID, email, password)
	if err != nil {
		log.Fatalf("Failed to create login: %v", err)
	}

	log.Printf("Admin user created successfully: %s (%s)", name, email)
}
