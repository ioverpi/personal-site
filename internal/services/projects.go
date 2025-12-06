package services

import (
	"database/sql"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
	"github.com/lib/pq"
)

type ProjectsService struct {
	app *app.App
}

func NewProjectsService(app *app.App) *ProjectsService {
	return &ProjectsService{app: app}
}

func (s *ProjectsService) GetAllProjects() ([]models.Project, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, name, description, tags, github_url, demo_url, display_order, created_at
		FROM projects
		ORDER BY display_order ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjects(rows)
}

func (s *ProjectsService) GetProjectByID(id int) (*models.Project, error) {
	var project models.Project
	err := s.app.DB.QueryRow(`
		SELECT id, name, description, tags, github_url, demo_url, display_order, created_at
		FROM projects
		WHERE id = $1
	`, id).Scan(
		&project.ID, &project.Name, &project.Description,
		pq.Array(&project.Tags), &project.GithubURL, &project.DemoURL,
		&project.DisplayOrder, &project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func scanProjects(rows *sql.Rows) ([]models.Project, error) {
	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description,
			pq.Array(&project.Tags), &project.GithubURL, &project.DemoURL,
			&project.DisplayOrder, &project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}
