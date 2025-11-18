package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	DB      string `json:"db,omitempty"`
}

// GetHealth handles GET /api/health
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Message: "Service is healthy",
	}

	// Check database connectivity
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			response.Status = "degraded"
			response.Message = "Service is running but database is unavailable"
			response.DB = "unavailable"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(response)
			return
		}
		response.DB = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
