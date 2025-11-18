package models

import "time"

// TaskDependency represents a dependency relationship between tasks
type TaskDependency struct {
	ID              int       `json:"id"`
	TaskID          int       `json:"task_id"`
	DependsOnTaskID int       `json:"depends_on_task_id"`
	CreatedAt       time.Time `json:"created_at"`
}
