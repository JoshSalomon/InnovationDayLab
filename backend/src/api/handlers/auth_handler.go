package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"

	"task-management-app/api/middleware"
	"task-management-app/api/response"
	"task-management-app/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
	store       *sessions.CookieStore
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db),
		store:       store,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handles POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		response.WriteError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create session
	session, err := h.store.Get(r, "session")
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	session.Values["userID"] = user.ID
	session.Values["userType"] = user.UserType
	session.Values["username"] = user.Username

	if err := session.Save(r, w); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to save session")
		return
	}

	// Return user info (without password hash)
	userResponse := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"display_name": user.DisplayName,
		"user_type":   user.UserType,
		"created_at":  user.CreatedAt,
	}

	response.WriteJSON(w, http.StatusOK, userResponse)
}

// Logout handles POST /api/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "session")
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to save session")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetCurrentUser handles GET /api/me (get current user info)
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userType, _ := middleware.GetUserType(r.Context())
	username, _ := r.Context().Value("username").(string)

	userResponse := map[string]interface{}{
		"id":       userID,
		"username": username,
		"user_type": userType,
	}

	response.WriteJSON(w, http.StatusOK, userResponse)
}
