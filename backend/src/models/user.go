package models

import "time"

// User represents an application user
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never serialize password hash
	Email        string    `json:"email"`
	DisplayName  *string   `json:"display_name,omitempty"`
	UserType     string    `json:"user_type"` // "admin" or "regular"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsAdmin returns true if the user is an admin
func (u *User) IsAdmin() bool {
	return u.UserType == "admin"
}
