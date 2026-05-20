package entity

import "testing"

func TestNewTaskDefaults(t *testing.T) {
	description := "write tests"

	task := NewTask(description)
	if task == nil {
		t.Fatal("NewTask() returned nil")
	}

	if task.ID != 1 {
		t.Fatalf("task.ID = %d, want 1", task.ID)
	}
	if task.Description != description {
		t.Fatalf("task.Description = %q, want %q", task.Description, description)
	}
	if task.Status != TaskStatusToDo {
		t.Fatalf("task.Status = %q, want %q", task.Status, TaskStatusToDo)
	}
	if task.CreatedAt.IsZero() {
		t.Fatal("task.CreatedAt should not be zero")
	}
	if task.UpdatedAt.IsZero() {
		t.Fatal("task.UpdatedAt should not be zero")
	}
	if task.UpdatedAt.Before(task.CreatedAt) {
		t.Fatalf("task.UpdatedAt (%v) is before task.CreatedAt (%v)", task.UpdatedAt, task.CreatedAt)
	}
}
