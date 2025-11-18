package services

import (
	"database/sql"
	"errors"
	"strings"

	"task-management-app/models"
)

// UserService handles user-related business logic
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	Username    string
	Password    string
	Email       string
	DisplayName *string
	UserType    string // "admin" or "regular"
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(req CreateUserRequest) (*models.User, error) {
	// Validate required fields
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("username is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("password is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email is required")
	}

	// Validate username length
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}

	// Validate password length
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Validate email format (basic check)
	if !strings.Contains(req.Email, "@") {
		return nil, errors.New("invalid email format")
	}

	// Validate display name length if provided
	if req.DisplayName != nil && len(*req.DisplayName) > 100 {
		return nil, errors.New("display name must be 100 characters or less")
	}

	// Validate user type
	userType := strings.ToLower(strings.TrimSpace(req.UserType))
	if userType == "" {
		userType = "regular" // Default to regular user
	}
	if userType != "admin" && userType != "regular" {
		return nil, errors.New("user_type must be 'admin' or 'regular'")
	}

	// Hash password
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Check if username already exists
	var existingUsername string
	err = s.db.QueryRow("SELECT username FROM users WHERE username = ?", req.Username).Scan(&existingUsername)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, errors.New("username already exists")
		}
		return nil, err
	}

	// Check if email already exists
	var existingEmail string
	err = s.db.QueryRow("SELECT email FROM users WHERE email = ?", req.Email).Scan(&existingEmail)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	// Insert new user
	result, err := s.db.Exec(`
		INSERT INTO users (username, password_hash, email, display_name, user_type)
		VALUES (?, ?, ?, ?, ?)
	`, req.Username, passwordHash, req.Email, req.DisplayName, userType)
	if err != nil {
		return nil, err
	}

	// Get the created user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the created user
	var user models.User
	var displayName sql.NullString
	err = s.db.QueryRow(`
		SELECT id, username, password_hash, email, display_name, user_type, created_at, updated_at
		FROM users
		WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&displayName,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if displayName.Valid {
		user.DisplayName = &displayName.String
	}

	return &user, nil
}
