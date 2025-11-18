package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"task-management-app/api/middleware"
	"task-management-app/api/response"
	"task-management-app/services"
)

// TaskHandler handles task-related endpoints
type TaskHandler struct {
	taskService       *services.TaskService
	dependencyService *services.DependencyService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{
		taskService:       services.NewTaskService(db),
		dependencyService: services.NewDependencyService(db),
	}
}

// GetTasks handles GET /api/tasks
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Get status filter from query parameter
	statusFilter := r.URL.Query().Get("status")
	if statusFilter == "" {
		statusFilter = "all"
	}

	// Validate status filter
	validFilters := map[string]bool{
		"all":         true,
		"completed":   true,
		"in_progress": true,
	}
	if !validFilters[statusFilter] {
		response.WriteError(w, http.StatusBadRequest, "Invalid status filter. Must be 'all', 'completed', or 'in_progress'")
		return
	}

	// Get tasks for user
	tasks, err := h.taskService.GetTasksByUserID(userID, statusFilter)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	// Convert tasks to response format
	taskResponses := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		// Get dependencies for this task
		dependencyIDs, _ := h.dependencyService.GetDependencyIDs(task.ID, userID)

		taskResponse := map[string]interface{}{
			"id":         task.ID,
			"user_id":    task.UserID,
			"description": task.Description,
			"status":     task.Status,
			"progress":   task.Progress,
			"created_at": task.CreatedAt,
			"updated_at": task.UpdatedAt,
			"dependencies": dependencyIDs,
		}
		if task.DueDate != nil {
			taskResponse["due_date"] = task.DueDate.Format("2006-01-02")
		}
		taskResponses[i] = taskResponse
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"tasks": taskResponses,
	})
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Description string  `json:"description"`
	DueDate     *string `json:"due_date,omitempty"` // ISO date string (YYYY-MM-DD)
	Status      string  `json:"status"`
	Progress    int     `json:"progress"`
}

// CreateTask handles POST /api/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		dueDate = &parsedDate
	}

	// Create task via service
	task, err := h.taskService.CreateTask(userID, req.Description, dueDate, req.Status, req.Progress)
	if err != nil {
		// Check for validation errors
		if err.Error() == "description is required" || err.Error() == "invalid status" || err.Error() == "progress must be between 0 and 100" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Get dependencies for this task (empty for new task)
	dependencyIDs, _ := h.dependencyService.GetDependencyIDs(task.ID, userID)

	// Return task response
	taskResponse := map[string]interface{}{
		"id":         task.ID,
		"user_id":    task.UserID,
		"description": task.Description,
		"status":     task.Status,
		"progress":   task.Progress,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
		"dependencies": dependencyIDs,
	}
	if task.DueDate != nil {
		taskResponse["due_date"] = task.DueDate.Format("2006-01-02")
	}

	response.WriteJSON(w, http.StatusCreated, taskResponse)
}

// GetTask handles GET /api/tasks/{taskId}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task ID from URL (chi router provides this via URLParam)
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task (ensuring it belongs to user)
	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve task")
		return
	}

	// Get dependencies for this task
	dependencyIDs, _ := h.dependencyService.GetDependencyIDs(task.ID, userID)

	// Return task response
	taskResponse := map[string]interface{}{
		"id":         task.ID,
		"user_id":    task.UserID,
		"description": task.Description,
		"status":     task.Status,
		"progress":   task.Progress,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
		"dependencies": dependencyIDs,
	}
	if task.DueDate != nil {
		taskResponse["due_date"] = task.DueDate.Format("2006-01-02")
	}

	response.WriteJSON(w, http.StatusOK, taskResponse)
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"` // ISO date string (YYYY-MM-DD) or null to clear
	Status      *string `json:"status,omitempty"`
	Progress    *int    `json:"progress,omitempty"`
}

// UpdateTask handles PUT /api/tasks/{taskId}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task ID from URL (chi router provides this via URLParam)
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != nil {
		if *req.DueDate == "" {
			// Empty string means clear the due date
			dueDate = nil
		} else {
			parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
				return
			}
			dueDate = &parsedDate
		}
	}

	// Update task via service
	task, err := h.taskService.UpdateTask(taskID, userID, services.UpdateTaskRequest{
		Description: req.Description,
		DueDate:     dueDate,
		Status:      req.Status,
		Progress:    req.Progress,
	})
	if err != nil {
		// Check for validation errors (should return 400)
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "cannot mark task as completed") ||
			strings.Contains(errorMsg, "description cannot be empty") ||
			strings.Contains(errorMsg, "invalid status") ||
			strings.Contains(errorMsg, "progress must be between") {
			response.WriteError(w, http.StatusBadRequest, errorMsg)
			return
		}
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	// Return updated task response
	taskResponse := map[string]interface{}{
		"id":         task.ID,
		"user_id":    task.UserID,
		"description": task.Description,
		"status":     task.Status,
		"progress":   task.Progress,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
	}
	if task.DueDate != nil {
		taskResponse["due_date"] = task.DueDate.Format("2006-01-02")
	}

	response.WriteJSON(w, http.StatusOK, taskResponse)
}

// DeleteTask handles DELETE /api/tasks/{taskId}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse task ID from URL (chi router provides this via URLParam)
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Delete task via service
	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
