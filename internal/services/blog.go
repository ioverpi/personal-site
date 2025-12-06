package services

import (
	"database/sql"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
	"github.com/lib/pq"
)

type BlogService struct {
	app *app.App
}

func NewBlogService(app *app.App) *BlogService {
	return &BlogService{app: app}
}

func (s *BlogService) GetPublishedPosts() ([]models.Post, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, title, slug, content, published_at, created_at, updated_at
		FROM posts
		WHERE published_at IS NOT NULL
		ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

func (s *BlogService) GetAllPosts() ([]models.Post, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, title, slug, content, published_at, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

func (s *BlogService) GetPostBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := s.app.DB.QueryRow(`
		SELECT id, title, slug, content, published_at, created_at, updated_at
		FROM posts
		WHERE slug = $1
	`, slug).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content,
		&post.PublishedAt, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *BlogService) GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	err := s.app.DB.QueryRow(`
		SELECT id, title, slug, content, published_at, created_at, updated_at
		FROM posts
		WHERE id = $1
	`, id).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content,
		&post.PublishedAt, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content,
			&post.PublishedAt, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// Ensure pq is imported for array handling
var _ = pq.Array
