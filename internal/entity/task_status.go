package entity

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
