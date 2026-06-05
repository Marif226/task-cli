package command

import (
	"errors"
	"testing"

	"github.com/Marif226/task-cli/internal/entity"
)

type fakeTaskStatusUpdater struct {
	calledID     int
	calledStatus entity.TaskStatus
	task         *entity.Task
	err          error
}

func (f *fakeTaskStatusUpdater) UpdateStatus(id int, status entity.TaskStatus) (*entity.Task, error) {
	f.calledID = id
	f.calledStatus = status
	if f.err != nil {
		return nil, f.err
	}
	if f.task != nil {
		return f.task, nil
	}
	return &entity.Task{ID: id, Status: status}, nil
}

func TestMarkInit(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantError  bool
		wantID     int
		wantStatus entity.TaskStatus
	}{
		{name: "missing args", args: nil, wantError: true},
		{name: "only id", args: []string{"1"}, wantError: true},
		{name: "invalid id", args: []string{"x", "done"}, wantError: true},
		{name: "invalid status", args: []string{"1", "blocked"}, wantError: true},
		{name: "valid args", args: []string{"4", "in-progress"}, wantID: 4, wantStatus: entity.TaskStatusInProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewMarkCommand(&fakeTaskStatusUpdater{})
			err := cmd.Init(tt.args)

			if tt.wantError && err == nil {
				t.Fatal("Init() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
			if !tt.wantError {
				if cmd.id != tt.wantID {
					t.Fatalf("id = %d, want %d", cmd.id, tt.wantID)
				}
				if cmd.status != tt.wantStatus {
					t.Fatalf("status = %q, want %q", cmd.status, tt.wantStatus)
				}
			}
		})
	}
}

func TestMarkRun(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		status     entity.TaskStatus
		storageErr error
		wantError  bool
	}{
		{name: "invalid id", id: 0, status: entity.TaskStatusDone, wantError: true},
		{name: "success forwards parsed status", id: 2, status: entity.TaskStatusDone, wantError: false},
		{name: "storage error propagated", id: 2, status: entity.TaskStatusToDo, storageErr: errors.New("write failed"), wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeTaskStatusUpdater{err: tt.storageErr}
			cmd := NewMarkCommand(storage)
			cmd.id = tt.id
			cmd.status = tt.status

			err := cmd.Run()
			if tt.wantError && err == nil {
				t.Fatal("Run() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}

			if tt.id > 0 {
				if storage.calledID != tt.id {
					t.Fatalf("UpdateStatus() id = %d, want %d", storage.calledID, tt.id)
				}
				if storage.calledStatus != tt.status {
					t.Fatalf("UpdateStatus() status = %q, want %q", storage.calledStatus, tt.status)
				}
			}
		})
	}
}
