package services

import (
	"database/sql"
	"errors"

	"task-management-app/models"
)

// DependencyService handles task dependency-related business logic
type DependencyService struct {
	db *sql.DB
}

// NewDependencyService creates a new dependency service
func NewDependencyService(db *sql.DB) *DependencyService {
	return &DependencyService{db: db}
}

// AddDependency adds a dependency relationship between two tasks
// It validates ownership, prevents self-dependency, and detects circular dependencies
func (s *DependencyService) AddDependency(taskID, dependsOnTaskID, userID int) error {
	// Prevent self-dependency
	if taskID == dependsOnTaskID {
		return errors.New("task cannot depend on itself")
	}

	// Verify both tasks exist and belong to the user
	task1, err := s.getTaskByIDAndUser(taskID, userID)
	if err != nil {
		return errors.New("task not found")
	}

	task2, err := s.getTaskByIDAndUser(dependsOnTaskID, userID)
	if err != nil {
		return errors.New("depends_on_task not found")
	}

	// Ensure both tasks belong to the same user
	if task1.UserID != userID || task2.UserID != userID {
		return errors.New("tasks must belong to the same user")
	}

	// Check if dependency already exists
	var exists int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM task_dependencies
		WHERE task_id = ? AND depends_on_task_id = ?
	`, taskID, dependsOnTaskID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("dependency already exists")
	}

	// Check for circular dependency (would adding this create a cycle?)
	if err := s.checkCircularDependency(taskID, dependsOnTaskID, userID); err != nil {
		return err
	}

	// Add dependency
	_, err = s.db.Exec(`
		INSERT INTO task_dependencies (task_id, depends_on_task_id)
		VALUES (?, ?)
	`, taskID, dependsOnTaskID)
	if err != nil {
		return err
	}

	return nil
}

// GetDependencies retrieves all tasks that a given task depends on
func (s *DependencyService) GetDependencies(taskID, userID int) ([]*models.Task, error) {
	// Verify task exists and belongs to user
	_, err := s.getTaskByIDAndUser(taskID, userID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// Get all dependencies
	rows, err := s.db.Query(`
		SELECT t.id, t.user_id, t.description, t.due_date, t.status, t.progress, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN task_dependencies td ON t.id = td.depends_on_task_id
		WHERE td.task_id = ? AND t.user_id = ?
		ORDER BY t.created_at DESC
	`, taskID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dependencies []*models.Task
	for rows.Next() {
		task, err := models.ScanTask(rows)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}

// RemoveDependency removes a dependency relationship
func (s *DependencyService) RemoveDependency(taskID, dependsOnTaskID, userID int) error {
	// Verify task exists and belongs to user
	_, err := s.getTaskByIDAndUser(taskID, userID)
	if err != nil {
		return errors.New("task not found")
	}

	// Remove dependency
	result, err := s.db.Exec(`
		DELETE FROM task_dependencies
		WHERE task_id = ? AND depends_on_task_id = ?
	`, taskID, dependsOnTaskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("dependency not found")
	}

	return nil
}

// GetDependencyIDs returns a list of task IDs that the given task depends on
func (s *DependencyService) GetDependencyIDs(taskID, userID int) ([]int, error) {
	dependencies, err := s.GetDependencies(taskID, userID)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(dependencies))
	for i, dep := range dependencies {
		ids[i] = dep.ID
	}

	return ids, nil
}

// checkCircularDependency checks if adding a dependency would create a circular dependency
// Uses DFS to detect cycles in the dependency graph
func (s *DependencyService) checkCircularDependency(taskID, dependsOnTaskID, userID int) error {
	// If taskID depends on dependsOnTaskID, check if dependsOnTaskID (or any of its dependencies) depends on taskID
	visited := make(map[int]bool)
	return s.dfsCheckCycle(dependsOnTaskID, taskID, userID, visited)
}

// dfsCheckCycle performs depth-first search to detect if there's a path from start to target
func (s *DependencyService) dfsCheckCycle(start, target, userID int, visited map[int]bool) error {
	if start == target {
		return errors.New("circular dependency detected")
	}

	if visited[start] {
		return nil // Already visited this node, no cycle through this path
	}

	visited[start] = true

	// Get all tasks that 'start' depends on
	rows, err := s.db.Query(`
		SELECT depends_on_task_id
		FROM task_dependencies
		WHERE task_id = ?
	`, start)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dependsOnID int
		if err := rows.Scan(&dependsOnID); err != nil {
			return err
		}

		// Verify the dependency belongs to the same user
		var depUserID int
		err := s.db.QueryRow("SELECT user_id FROM tasks WHERE id = ?", dependsOnID).Scan(&depUserID)
		if err != nil {
			continue // Skip if task doesn't exist
		}
		if depUserID != userID {
			continue // Skip if not same user
		}

		// Recursively check if this dependency leads to target
		if err := s.dfsCheckCycle(dependsOnID, target, userID, visited); err != nil {
			return err
		}
	}

	return rows.Err()
}

// getTaskByIDAndUser is a helper to verify task ownership
func (s *DependencyService) getTaskByIDAndUser(taskID, userID int) (*models.Task, error) {
	var task models.Task
	var dueDate sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, user_id, description, due_date, status, progress, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Description,
		&dueDate,
		&task.Status,
		&task.Progress,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}

	return &task, nil
}
