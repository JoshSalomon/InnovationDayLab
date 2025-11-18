package services

import (
	"database/sql"
	"errors"

	"task-management-app/models"
)

// AuthService handles authentication operations
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// Login authenticates a user by username and password
func (s *AuthService) Login(username, password string) (*models.User, error) {
	var user models.User
	var passwordHash string

	err := s.db.QueryRow(`
		SELECT id, username, password_hash, email, display_name, user_type, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&passwordHash,
		&user.Email,
		&user.DisplayName,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid credentials")
	}
	if err != nil {
		return nil, err
	}

	// Verify password
	if !CheckPassword(password, passwordHash) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
