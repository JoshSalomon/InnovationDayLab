package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"task-management-app/api/middleware"
	"task-management-app/api/response"
	"task-management-app/services"
)

// DependencyHandler handles dependency-related endpoints
type DependencyHandler struct {
	dependencyService *services.DependencyService
}

// NewDependencyHandler creates a new dependency handler
func NewDependencyHandler(db *sql.DB) *DependencyHandler {
	return &DependencyHandler{
		dependencyService: services.NewDependencyService(db),
	}
}

// GetDependencies handles GET /api/tasks/{taskId}/dependencies
func (h *DependencyHandler) GetDependencies(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task ID from URL
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get dependencies
	dependencies, err := h.dependencyService.GetDependencies(taskID, userID)
	if err != nil {
		if err.Error() == "task not found" {
			response.WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve dependencies")
		return
	}

	// Convert to response format
	dependencyResponses := make([]map[string]interface{}, len(dependencies))
	for i, dep := range dependencies {
		depResp := map[string]interface{}{
			"id":         dep.ID,
			"user_id":    dep.UserID,
			"description": dep.Description,
			"status":     dep.Status,
			"progress":   dep.Progress,
			"created_at": dep.CreatedAt,
			"updated_at": dep.UpdatedAt,
		}
		if dep.DueDate != nil {
			depResp["due_date"] = dep.DueDate.Format("2006-01-02")
		}
		dependencyResponses[i] = depResp
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dependencies": dependencyResponses,
	})
}

// AddDependencyRequest represents the request body for adding a dependency
type AddDependencyRequest struct {
	DependsOnTaskID int `json:"depends_on_task_id"`
}

// AddDependency handles POST /api/tasks/{taskId}/dependencies
func (h *DependencyHandler) AddDependency(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task ID from URL
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Parse request body
	var req AddDependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate depends_on_task_id
	if req.DependsOnTaskID <= 0 {
		response.WriteError(w, http.StatusBadRequest, "Invalid depends_on_task_id")
		return
	}

	// Add dependency
	err = h.dependencyService.AddDependency(taskID, req.DependsOnTaskID, userID)
	if err != nil {
		if err.Error() == "task not found" || err.Error() == "depends_on_task not found" {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "task cannot depend on itself" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "dependency already exists" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "circular dependency detected" {
			response.WriteError(w, http.StatusBadRequest, "Circular dependency detected: This would create a circular dependency")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to add dependency")
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Dependency added successfully",
	})
}

// RemoveDependency handles DELETE /api/tasks/{taskId}/dependencies/{dependsOnTaskId}
func (h *DependencyHandler) RemoveDependency(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task IDs from URL
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	dependsOnTaskIDStr := chi.URLParam(r, "dependsOnTaskId")
	dependsOnTaskID, err := strconv.Atoi(dependsOnTaskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid depends_on_task_id")
		return
	}

	// Remove dependency
	err = h.dependencyService.RemoveDependency(taskID, dependsOnTaskID, userID)
	if err != nil {
		if err.Error() == "task not found" || err.Error() == "dependency not found" {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to remove dependency")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
