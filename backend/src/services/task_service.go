package services

import (
	"database/sql"
	"errors"
	"time"

	"task-management-app/models"
)

// TaskService handles task-related business logic
type TaskService struct {
	db               *sql.DB
	dependencyService *DependencyService
}

// NewTaskService creates a new task service
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{
		db:               db,
		dependencyService: NewDependencyService(db),
	}
}

// GetTasksByUserID retrieves all tasks for a specific user, optionally filtered by status
func (s *TaskService) GetTasksByUserID(userID int, statusFilter string) ([]*models.Task, error) {
	// Ensure data isolation - always filter by user_id (FR-008)
	query := `
		SELECT id, user_id, description, due_date, status, progress, created_at, updated_at
		FROM tasks
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	// Apply status filter if specified
	if statusFilter != "" && statusFilter != "all" {
		if statusFilter == "completed" {
			query += " AND status = 'completed'"
		} else if statusFilter == "in_progress" {
			// "in_progress" means all non-completed and non-deferred tasks
			query += " AND status != 'completed' AND status != 'deferred'"
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task, err := models.ScanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// CreateTask creates a new task for a user
func (s *TaskService) CreateTask(userID int, description string, dueDate *time.Time, status string, progress int) (*models.Task, error) {
	// Validate required fields
	if description == "" {
		return nil, errors.New("description is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"completed":   true,
		"in_progress": true,
		"not_started": true,
		"blocked":     true,
		"deferred":    true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	// Validate progress (0-100%)
	if progress < 0 || progress > 100 {
		return nil, errors.New("progress must be between 0 and 100")
	}

	// Insert task
	result, err := s.db.Exec(`
		INSERT INTO tasks (user_id, description, due_date, status, progress)
		VALUES (?, ?, ?, ?, ?)
	`, userID, description, dueDate, status, progress)
	if err != nil {
		return nil, err
	}

	// Get the created task ID
	taskID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the created task
	var task models.Task
	var dueDateNull sql.NullTime
	err = s.db.QueryRow(`
		SELECT id, user_id, description, due_date, status, progress, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`, taskID).Scan(
		&task.ID,
		&task.UserID,
		&task.Description,
		&dueDateNull,
		&task.Status,
		&task.Progress,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if dueDateNull.Valid {
		task.DueDate = &dueDateNull.Time
	}

	return &task, nil
}

// GetTaskByID retrieves a specific task by ID, ensuring it belongs to the user
func (s *TaskService) GetTaskByID(taskID, userID int) (*models.Task, error) {
	var task models.Task
	var dueDateNull sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, user_id, description, due_date, status, progress, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Description,
		&dueDateNull,
		&task.Status,
		&task.Progress,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if dueDateNull.Valid {
		task.DueDate = &dueDateNull.Time
	}

	return &task, nil
}

// UpdateTaskRequest represents a task update request
type UpdateTaskRequest struct {
	Description *string
	DueDate     *time.Time
	Status      *string
	Progress    *int
}

// UpdateTask updates an existing task, ensuring it belongs to the user
func (s *TaskService) UpdateTask(taskID, userID int, req UpdateTaskRequest) (*models.Task, error) {
	// First, verify the task exists and belongs to the user
	existingTask, err := s.GetTaskByID(taskID, userID)
	if err != nil {
		return nil, err
	}

	// Build update query dynamically based on provided fields
	updates := []string{}
	args := []interface{}{}

	if req.Description != nil {
		if *req.Description == "" {
			return nil, errors.New("description cannot be empty")
		}
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
	}

	if req.DueDate != nil {
		updates = append(updates, "due_date = ?")
		args = append(args, *req.DueDate)
	}

	if req.Status != nil {
		// Validate status
		validStatuses := map[string]bool{
			"completed":   true,
			"in_progress": true,
			"not_started": true,
			"blocked":     true,
			"deferred":    true,
		}
		if !validStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
		
		// If marking task as completed, check that all dependencies are completed
		if *req.Status == "completed" {
			dependencies, err := s.dependencyService.GetDependencies(taskID, userID)
			if err != nil {
				return nil, err
			}
			
			// Check if any dependency is not completed
			var incompleteDeps []string
			for _, dep := range dependencies {
				if dep.Status != "completed" {
					incompleteDeps = append(incompleteDeps, dep.Description)
				}
			}
			
			if len(incompleteDeps) > 0 {
				return nil, errors.New("cannot mark task as completed: one or more dependencies are not completed")
			}
		}
		
		updates = append(updates, "status = ?")
		args = append(args, *req.Status)
	}

	if req.Progress != nil {
		// Validate progress (0-100%)
		if *req.Progress < 0 || *req.Progress > 100 {
			return nil, errors.New("progress must be between 0 and 100")
		}
		updates = append(updates, "progress = ?")
		args = append(args, *req.Progress)
	}

	if len(updates) == 0 {
		// No updates provided, return existing task
		return existingTask, nil
	}

	// Add updated_at
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, taskID, userID)

	// Execute update
	query := "UPDATE tasks SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ? AND user_id = ?"

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Retrieve updated task
	return s.GetTaskByID(taskID, userID)
}

// DeleteTask deletes a task, ensuring it belongs to the user
// Dependencies are automatically removed by database CASCADE DELETE
func (s *TaskService) DeleteTask(taskID, userID int) error {
	// First, verify the task exists and belongs to the user
	_, err := s.GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}

	// Delete the task (CASCADE DELETE will handle dependencies automatically)
	_, err = s.db.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", taskID, userID)
	if err != nil {
		return err
	}

	return nil
}
