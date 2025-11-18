package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"task-management-app/api/middleware"
	"task-management-app/api/response"
	"task-management-app/services"
)

// UserHandler handles user management endpoints
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(db),
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name,omitempty"`
	UserType    string  `json:"user_type,omitempty"` // Optional, defaults to "regular"
}

// CreateUser handles POST /api/users (admin only)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin (admin middleware should be applied, but double-check)
	userType, ok := middleware.GetUserType(r.Context())
	if !ok || userType != "admin" {
		response.WriteError(w, http.StatusForbidden, "Only admin users can create new users")
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set default user type to "regular" if not specified
	if req.UserType == "" {
		req.UserType = "regular"
	}

	// Create user via service
	user, err := h.userService.CreateUser(services.CreateUserRequest{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		UserType:    req.UserType,
	})
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if errMsg == "username already exists" || errMsg == "email already exists" {
			response.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}
		if strings.Contains(errMsg, "required") || strings.Contains(errMsg, "must be") {
			response.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return user response (without password hash)
	userResponse := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"display_name": user.DisplayName,
		"user_type":   user.UserType,
		"created_at":  user.CreatedAt,
		"updated_at":  user.UpdatedAt,
	}

	response.WriteJSON(w, http.StatusCreated, userResponse)
}
