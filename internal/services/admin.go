package services

import (
	"strings"
	"time"
	"unicode"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
	"github.com/lib/pq"
)

type AdminService struct {
	app *app.App
}

func NewAdminService(app *app.App) *AdminService {
	return &AdminService{app: app}
}

// Posts

type CreatePostInput struct {
	Title   string
	Slug    string
	Content string
	Publish bool
}

type UpdatePostInput struct {
	Title   string
	Slug    string
	Content string
	Publish bool
}

func (s *AdminService) CreatePost(input CreatePostInput) (*models.Post, error) {
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Title)
	}

	var publishedAt *time.Time
	if input.Publish {
		now := time.Now()
		publishedAt = &now
	}

	var post models.Post
	err := s.app.DB.QueryRow(`
		INSERT INTO posts (title, slug, content, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, title, slug, content, published_at, created_at, updated_at
	`, input.Title, slug, input.Content, publishedAt).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content,
		&post.PublishedAt, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *AdminService) UpdatePost(id int, input UpdatePostInput) (*models.Post, error) {
	// Get current post to check publish status
	var currentPublishedAt *time.Time
	err := s.app.DB.QueryRow(`SELECT published_at FROM posts WHERE id = $1`, id).Scan(&currentPublishedAt)
	if err != nil {
		return nil, err
	}

	var publishedAt *time.Time
	if input.Publish {
		if currentPublishedAt != nil {
			publishedAt = currentPublishedAt // Keep original publish date
		} else {
			now := time.Now()
			publishedAt = &now
		}
	}

	var post models.Post
	err = s.app.DB.QueryRow(`
		UPDATE posts
		SET title = $1, slug = $2, content = $3, published_at = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, slug, content, published_at, created_at, updated_at
	`, input.Title, input.Slug, input.Content, publishedAt, id).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content,
		&post.PublishedAt, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *AdminService) DeletePost(id int) error {
	_, err := s.app.DB.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}

// Projects

type CreateProjectInput struct {
	Name         string
	Description  string
	Tags         []string
	GithubURL    string
	DemoURL      string
	DisplayOrder int
}

type UpdateProjectInput struct {
	Name         string
	Description  string
	Tags         []string
	GithubURL    string
	DemoURL      string
	DisplayOrder int
}

func (s *AdminService) CreateProject(input CreateProjectInput) (*models.Project, error) {
	var githubURL, demoURL *string
	if input.GithubURL != "" {
		githubURL = &input.GithubURL
	}
	if input.DemoURL != "" {
		demoURL = &input.DemoURL
	}

	var project models.Project
	err := s.app.DB.QueryRow(`
		INSERT INTO projects (name, description, tags, github_url, demo_url, display_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, name, description, tags, github_url, demo_url, display_order, created_at
	`, input.Name, input.Description, pq.Array(input.Tags), githubURL, demoURL, input.DisplayOrder).Scan(
		&project.ID, &project.Name, &project.Description,
		pq.Array(&project.Tags), &project.GithubURL, &project.DemoURL,
		&project.DisplayOrder, &project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *AdminService) UpdateProject(id int, input UpdateProjectInput) (*models.Project, error) {
	var githubURL, demoURL *string
	if input.GithubURL != "" {
		githubURL = &input.GithubURL
	}
	if input.DemoURL != "" {
		demoURL = &input.DemoURL
	}

	var project models.Project
	err := s.app.DB.QueryRow(`
		UPDATE projects
		SET name = $1, description = $2, tags = $3, github_url = $4, demo_url = $5, display_order = $6
		WHERE id = $7
		RETURNING id, name, description, tags, github_url, demo_url, display_order, created_at
	`, input.Name, input.Description, pq.Array(input.Tags), githubURL, demoURL, input.DisplayOrder, id).Scan(
		&project.ID, &project.Name, &project.Description,
		pq.Array(&project.Tags), &project.GithubURL, &project.DemoURL,
		&project.DisplayOrder, &project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *AdminService) DeleteProject(id int) error {
	_, err := s.app.DB.Exec(`DELETE FROM projects WHERE id = $1`, id)
	return err
}

// Quotes

type CreateQuoteInput struct {
	Content string
	Author  string
	IsOwn   bool
}

type UpdateQuoteInput struct {
	Content string
	Author  string
	IsOwn   bool
}

func (s *AdminService) CreateQuote(input CreateQuoteInput) (*models.Quote, error) {
	var quote models.Quote
	err := s.app.DB.QueryRow(`
		INSERT INTO quotes (content, author, is_own, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, content, author, is_own, created_at
	`, input.Content, input.Author, input.IsOwn).Scan(
		&quote.ID, &quote.Content, &quote.Author, &quote.IsOwn, &quote.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (s *AdminService) UpdateQuote(id int, input UpdateQuoteInput) (*models.Quote, error) {
	var quote models.Quote
	err := s.app.DB.QueryRow(`
		UPDATE quotes
		SET content = $1, author = $2, is_own = $3
		WHERE id = $4
		RETURNING id, content, author, is_own, created_at
	`, input.Content, input.Author, input.IsOwn, id).Scan(
		&quote.ID, &quote.Content, &quote.Author, &quote.IsOwn, &quote.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (s *AdminService) DeleteQuote(id int) error {
	_, err := s.app.DB.Exec(`DELETE FROM quotes WHERE id = $1`, id)
	return err
}

// Helper functions

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	var result strings.Builder

	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	// Remove consecutive dashes
	s := result.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	return strings.Trim(s, "-")
}
