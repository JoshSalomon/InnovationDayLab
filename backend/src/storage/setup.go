package storage

import (
	"database/sql"
	"log"
	"os"

	"task-management-app/services"
)

// EnsureAdminUser creates the initial admin user if it doesn't exist
func EnsureAdminUser(db *sql.DB) error {
	// Check if admin user exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE user_type = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	// Get admin credentials from environment or use defaults
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "admin123" // Default password - should be changed in production
		log.Println("WARNING: Using default admin password. Set ADMIN_PASSWORD environment variable.")
	}

	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		email = "admin@example.com"
	}

	// Hash password
	passwordHash, err := services.HashPassword(password)
	if err != nil {
		return err
	}

	// Create admin user
	_, err = db.Exec(`
		INSERT INTO users (username, password_hash, email, user_type)
		VALUES (?, ?, ?, 'admin')
	`, username, passwordHash, email)

	if err != nil {
		return err
	}

	log.Printf("Admin user created: %s", username)
	return nil
}
