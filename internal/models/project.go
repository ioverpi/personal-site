package models

import "time"

type Project struct {
	ID           int
	Name         string
	Description  string
	Tags         []string
	GithubURL    *string
	DemoURL      *string
	DisplayOrder int
	CreatedAt    time.Time
}
