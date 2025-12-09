package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid or expired session")
	ErrInvalidInvite      = errors.New("invalid or expired invite")
	ErrInviteAlreadyUsed  = errors.New("invite already used")
)

type AuthService struct {
	app *app.App
}

func NewAuthService(app *app.App) *AuthService {
	return &AuthService{app: app}
}

// Password hashing

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Token generation

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Login methods

func (s *AuthService) GetLoginByEmail(email string) (*models.Login, error) {
	var login models.Login
	err := s.app.DB.QueryRow(`
		SELECT id, user_id, provider, provider_id, password_hash, created_at, updated_at
		FROM logins
		WHERE provider = $1 AND provider_id = $2
	`, models.ProviderPassword, email).Scan(
		&login.ID, &login.UserID, &login.Provider, &login.ProviderID,
		&login.PasswordHash, &login.CreatedAt, &login.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &login, nil
}

func (s *AuthService) CreatePasswordLogin(userID int, email, password string) (*models.Login, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	var login models.Login
	err = s.app.DB.QueryRow(`
		INSERT INTO logins (user_id, provider, provider_id, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, provider, provider_id, password_hash, created_at, updated_at
	`, userID, models.ProviderPassword, email, hash).Scan(
		&login.ID, &login.UserID, &login.Provider, &login.ProviderID,
		&login.PasswordHash, &login.CreatedAt, &login.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &login, nil
}

func (s *AuthService) UpdatePassword(userID int, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = s.app.DB.Exec(`
		UPDATE logins
		SET password_hash = $1, updated_at = NOW()
		WHERE user_id = $2 AND provider = $3
	`, hash, userID, models.ProviderPassword)
	return err
}

// Authentication

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	login, err := s.GetLoginByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if login.PasswordHash == nil || !CheckPassword(password, *login.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	var user models.User
	err = s.app.DB.QueryRow(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, login.UserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Sessions

func (s *AuthService) CreateSession(userID int, duration time.Duration) (*models.Session, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(duration)

	var session models.Session
	err = s.app.DB.QueryRow(`
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token, expires_at, created_at
	`, userID, token, expiresAt).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *AuthService) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := s.app.DB.QueryRow(`
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1
	`, token).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidSession
		}
		return nil, err
	}

	if session.IsExpired() {
		s.DeleteSession(token)
		return nil, ErrInvalidSession
	}

	return &session, nil
}

func (s *AuthService) GetUserBySession(token string) (*models.User, error) {
	session, err := s.GetSession(token)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.app.DB.QueryRow(`
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, session.UserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) DeleteSession(token string) error {
	_, err := s.app.DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (s *AuthService) DeleteUserSessions(userID int) error {
	_, err := s.app.DB.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (s *AuthService) CleanExpiredSessions() error {
	_, err := s.app.DB.Exec(`DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}

// Invites

func (s *AuthService) CreateInvite(email string, invitedBy int, duration time.Duration) (*models.Invite, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(duration)

	var invite models.Invite
	err = s.app.DB.QueryRow(`
		INSERT INTO invites (email, token, invited_by, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, token, invited_by, used_at, expires_at, created_at
	`, email, token, invitedBy, expiresAt).Scan(
		&invite.ID, &invite.Email, &invite.Token, &invite.InvitedBy,
		&invite.UsedAt, &invite.ExpiresAt, &invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (s *AuthService) GetInvite(token string) (*models.Invite, error) {
	var invite models.Invite
	err := s.app.DB.QueryRow(`
		SELECT id, email, token, invited_by, used_at, expires_at, created_at
		FROM invites
		WHERE token = $1
	`, token).Scan(
		&invite.ID, &invite.Email, &invite.Token, &invite.InvitedBy,
		&invite.UsedAt, &invite.ExpiresAt, &invite.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidInvite
		}
		return nil, err
	}

	if invite.IsExpired() {
		return nil, ErrInvalidInvite
	}

	if invite.IsUsed() {
		return nil, ErrInviteAlreadyUsed
	}

	return &invite, nil
}

func (s *AuthService) UseInvite(token string) error {
	result, err := s.app.DB.Exec(`
		UPDATE invites
		SET used_at = NOW()
		WHERE token = $1 AND used_at IS NULL AND expires_at > NOW()
	`, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrInvalidInvite
	}

	return nil
}

func (s *AuthService) GetPendingInvites() ([]models.Invite, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, email, token, invited_by, used_at, expires_at, created_at
		FROM invites
		WHERE used_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []models.Invite
	for rows.Next() {
		var invite models.Invite
		err := rows.Scan(
			&invite.ID, &invite.Email, &invite.Token, &invite.InvitedBy,
			&invite.UsedAt, &invite.ExpiresAt, &invite.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}
	return invites, rows.Err()
}

func (s *AuthService) DeleteInvite(id int) error {
	_, err := s.app.DB.Exec(`DELETE FROM invites WHERE id = $1`, id)
	return err
}
