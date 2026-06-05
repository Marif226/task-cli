package entity

import "fmt"

// TaskStatus describes the current state of a task.
type TaskStatus string

const (
	// TaskStatusToDo marks a task that has not been started.
	TaskStatusToDo TaskStatus = "to-do"
	// TaskStatusInProgress marks a task that is currently being worked on.
	TaskStatusInProgress = "in-progress"
	// TaskStatusDone marks a completed task.
	TaskStatusDone = "done"
)

// IsValid reports whether the status is one of the known task statuses.
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusToDo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// ParseTaskStatus validates and converts a string into TaskStatus.
func ParseTaskStatus(status string) (TaskStatus, error) {
	ts := TaskStatus(status)
	if !ts.IsValid() {
		return "", fmt.Errorf("invalid status %q: allowed values are %q, %q, %q", status, TaskStatusToDo, TaskStatusInProgress, TaskStatusDone)
	}

	return ts, nil
}
