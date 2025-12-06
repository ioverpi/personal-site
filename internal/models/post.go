package models

import "time"

type Post struct {
	ID          int
	Title       string
	Slug        string
	Content     string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
