package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Marif226/task-cli/internal/entity"
)

func TestTaskCreatePersistsAndIncrementsID(t *testing.T) {
	taskStorage, filePath := newTempTaskStorage(t, `{"last_id":0,"tasks":[]}`)

	id1, err := taskStorage.Create("first")
	if err != nil {
		t.Fatalf("Create() first returned error: %v", err)
	}
	id2, err := taskStorage.Create("second")
	if err != nil {
		t.Fatalf("Create() second returned error: %v", err)
	}

	if id1 != 1 || id2 != 2 {
		t.Fatalf("unexpected IDs: got (%d, %d), want (1, 2)", id1, id2)
	}

	saved := readTaskStorageFile(t, filePath)
	if saved.LastID != 2 {
		t.Fatalf("LastID = %d, want 2", saved.LastID)
	}
	if len(saved.Tasks) != 2 {
		t.Fatalf("len(Tasks) = %d, want 2", len(saved.Tasks))
	}
	if saved.Tasks[0].Status != entity.TaskStatusToDo || saved.Tasks[1].Status != entity.TaskStatusToDo {
		t.Fatal("new tasks should be persisted with to-do status")
	}
}

func TestTaskDeleteSuccessAndNotFound(t *testing.T) {
	taskStorage, filePath := newTempTaskStorage(t, "")

	now := time.Now().Add(-time.Minute)
	writeTaskStorageFile(t, filePath, &entity.TaskStorage{
		LastID: 2,
		Tasks: []entity.Task{
			{ID: 1, Description: "one", Status: entity.TaskStatusToDo, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Description: "two", Status: entity.TaskStatusDone, CreatedAt: now, UpdatedAt: now},
		},
	})

	if err := taskStorage.Delete(1); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	saved := readTaskStorageFile(t, filePath)
	if len(saved.Tasks) != 1 || saved.Tasks[0].ID != 2 {
		t.Fatalf("unexpected tasks after delete: %+v", saved.Tasks)
	}

	if err := taskStorage.Delete(999); err == nil {
		t.Fatal("Delete() expected not-found error, got nil")
	}
}

func TestTaskListFiltersByStatus(t *testing.T) {
	taskStorage, filePath := newTempTaskStorage(t, "")

	now := time.Now().Add(-time.Minute)
	writeTaskStorageFile(t, filePath, &entity.TaskStorage{
		LastID: 4,
		Tasks: []entity.Task{
			{ID: 1, Description: "a", Status: entity.TaskStatusToDo, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Description: "b", Status: entity.TaskStatusInProgress, CreatedAt: now, UpdatedAt: now},
			{ID: 3, Description: "c", Status: entity.TaskStatusDone, CreatedAt: now, UpdatedAt: now},
			{ID: 4, Description: "d", Status: entity.TaskStatusToDo, CreatedAt: now, UpdatedAt: now},
		},
	})

	tests := []struct {
		name       string
		status     entity.TaskStatus
		wantLength int
	}{
		{name: "all tasks", status: "", wantLength: 4},
		{name: "todo", status: entity.TaskStatusToDo, wantLength: 2},
		{name: "in-progress", status: entity.TaskStatusInProgress, wantLength: 1},
		{name: "done", status: entity.TaskStatusDone, wantLength: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := taskStorage.List(tt.status)
			if err != nil {
				t.Fatalf("List() returned error: %v", err)
			}
			if len(tasks) != tt.wantLength {
				t.Fatalf("len(List()) = %d, want %d", len(tasks), tt.wantLength)
			}
		})
	}
}

func TestTaskUpdateStatusSuccessAndNotFound(t *testing.T) {
	taskStorage, filePath := newTempTaskStorage(t, "")

	oldTime := time.Now().Add(-2 * time.Hour)
	writeTaskStorageFile(t, filePath, &entity.TaskStorage{
		LastID: 1,
		Tasks: []entity.Task{
			{ID: 1, Description: "one", Status: entity.TaskStatusToDo, CreatedAt: oldTime, UpdatedAt: oldTime},
		},
	})

	updatedTask, err := taskStorage.UpdateStatus(1, entity.TaskStatusDone)
	if err != nil {
		t.Fatalf("UpdateStatus() returned error: %v", err)
	}
	if updatedTask.Status != entity.TaskStatusDone {
		t.Fatalf("updated status = %q, want %q", updatedTask.Status, entity.TaskStatusDone)
	}
	if !updatedTask.UpdatedAt.After(oldTime) {
		t.Fatalf("updated timestamp %v should be after %v", updatedTask.UpdatedAt, oldTime)
	}

	saved := readTaskStorageFile(t, filePath)
	if len(saved.Tasks) != 1 || saved.Tasks[0].Status != entity.TaskStatusDone {
		t.Fatalf("unexpected persisted tasks after update: %+v", saved.Tasks)
	}

	_, err = taskStorage.UpdateStatus(999, entity.TaskStatusToDo)
	if err == nil {
		t.Fatal("UpdateStatus() expected not-found error, got nil")
	}
}

func TestTaskErrorPathsMalformedJSONAndMissingFile(t *testing.T) {
	operations := []struct {
		name string
		run  func(*Task) error
	}{
		{
			name: "create",
			run: func(s *Task) error {
				_, err := s.Create("x")
				return err
			},
		},
		{
			name: "delete",
			run: func(s *Task) error {
				return s.Delete(1)
			},
		},
		{
			name: "list",
			run: func(s *Task) error {
				_, err := s.List("")
				return err
			},
		},
		{
			name: "update status",
			run: func(s *Task) error {
				_, err := s.UpdateStatus(1, entity.TaskStatusDone)
				return err
			},
		},
	}

	for _, op := range operations {
		t.Run("malformed JSON "+op.name, func(t *testing.T) {
			taskStorage, filePath := newTempTaskStorage(t, "")
			if err := os.WriteFile(filePath, []byte("{"), 0644); err != nil {
				t.Fatalf("failed writing malformed JSON: %v", err)
			}

			err := op.run(taskStorage)
			if err == nil {
				t.Fatal("expected error for malformed JSON, got nil")
			}
		})

		t.Run("missing file "+op.name, func(t *testing.T) {
			dir := t.TempDir()
			missingPath := filepath.Join(dir, "missing-tasks.json")
			taskStorage := NewTaskStorage(missingPath)

			err := op.run(taskStorage)
			if err == nil {
				t.Fatal("expected error for missing file, got nil")
			}
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected os.ErrNotExist, got %v", err)
			}
		})
	}
}
