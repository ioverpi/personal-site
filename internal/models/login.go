package models

import "time"

type Login struct {
	ID           int
	UserID       int
	Provider     string
	ProviderID   string
	PasswordHash *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const (
	ProviderPassword = "password"
	ProviderGoogle   = "google"
	ProviderGitHub   = "github"
)
