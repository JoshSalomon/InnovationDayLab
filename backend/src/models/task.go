package models

import (
	"database/sql"
	"time"
)

// Task represents a work item
type Task struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"` // "completed", "in_progress", "not_started", "blocked", "deferred"
	Progress    int        `json:"progress"` // 0-100
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ScanTask scans a task from database rows
func ScanTask(rows *sql.Rows) (*Task, error) {
	task := &Task{}
	var dueDate sql.NullTime

	err := rows.Scan(
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

	return task, nil
}
