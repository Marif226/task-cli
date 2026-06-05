package entity

import (
	"time"
)

// TaskStorage is the JSON representation of all saved tasks.
type TaskStorage struct {
	LastID int    `json:"last_id"`
	Tasks  []Task `json:"tasks"`
}

// Task represents a single tracked task.
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
