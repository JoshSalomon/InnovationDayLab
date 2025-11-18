package models

import (
	"testing"
)

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		userType string
		want     bool
	}{
		{
			name:     "admin user",
			userType: "admin",
			want:     true,
		},
		{
			name:     "regular user",
			userType: "regular",
			want:     false,
		},
		{
			name:     "empty user type",
			userType: "",
			want:     false,
		},
		{
			name:     "invalid user type",
			userType: "invalid",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				UserType: tt.userType,
			}
			if got := u.IsAdmin(); got != tt.want {
				t.Errorf("User.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsAdmin_WithFullUser(t *testing.T) {
	adminUser := &User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		UserType: "admin",
	}

	regularUser := &User{
		ID:       2,
		Username: "user",
		Email:    "user@example.com",
		UserType: "regular",
	}

	if !adminUser.IsAdmin() {
		t.Error("Admin user should return true for IsAdmin()")
	}

	if regularUser.IsAdmin() {
		t.Error("Regular user should return false for IsAdmin()")
	}
}
