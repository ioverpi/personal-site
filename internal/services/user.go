package services

import (
	"database/sql"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
)

type UserService struct {
	app *app.App
}

func NewUserService(app *app.App) *UserService {
	return &UserService{app: app}
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	var user models.User
	err := s.app.DB.QueryRow(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.app.DB.QueryRow(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanUsers(rows)
}

func (s *UserService) GetAdmins() ([]models.User, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE role = 'admin'
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanUsers(rows)
}

type CreateUserInput struct {
	Email string
	Name  string
	Role  string
}

func (s *UserService) CreateUser(input CreateUserInput) (*models.User, error) {
	var user models.User
	err := s.app.DB.QueryRow(`
		INSERT INTO users (email, name, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, role, created_at, updated_at
	`, input.Email, input.Name, input.Role).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type UpdateUserInput struct {
	Name string
	Role string
}

func (s *UserService) UpdateUser(id int, input UpdateUserInput) (*models.User, error) {
	var user models.User
	err := s.app.DB.QueryRow(`
		UPDATE users
		SET name = $1, role = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, name, role, created_at, updated_at
	`, input.Name, input.Role, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(id int) error {
	_, err := s.app.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func scanUsers(rows *sql.Rows) ([]models.User, error) {
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Role,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
