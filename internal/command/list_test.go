package command

import (
	"errors"
	"testing"

	"github.com/Marif226/task-cli/internal/entity"
)

type fakeTaskLister struct {
	calledWith entity.TaskStatus
	tasks      []entity.Task
	err        error
}

func (f *fakeTaskLister) List(status entity.TaskStatus) ([]entity.Task, error) {
	f.calledWith = status
	if f.err != nil {
		return nil, f.err
	}
	return f.tasks, nil
}

func TestListInit(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantError  bool
		wantStatus entity.TaskStatus
	}{
		{name: "no status filter", args: nil, wantStatus: ""},
		{name: "valid status", args: []string{"done"}, wantStatus: entity.TaskStatusDone},
		{name: "invalid status", args: []string{"blocked"}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewListCommand(&fakeTaskLister{})
			err := cmd.Init(tt.args)

			if tt.wantError && err == nil {
				t.Fatal("Init() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
			if !tt.wantError && cmd.status != tt.wantStatus {
				t.Fatalf("status = %q, want %q", cmd.status, tt.wantStatus)
			}
		})
	}
}

func TestListRun(t *testing.T) {
	tests := []struct {
		name       string
		status     entity.TaskStatus
		storageErr error
		wantError  bool
	}{
		{name: "success", status: entity.TaskStatusInProgress, wantError: false},
		{name: "storage error propagated", status: entity.TaskStatusDone, storageErr: errors.New("cannot read"), wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeTaskLister{
				tasks: []entity.Task{{ID: 1, Description: "demo", Status: tt.status}},
				err:   tt.storageErr,
			}
			cmd := NewListCommand(storage)
			cmd.status = tt.status

			err := cmd.Run()
			if tt.wantError && err == nil {
				t.Fatal("Run() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}
			if storage.calledWith != tt.status {
				t.Fatalf("List() called with %q, want %q", storage.calledWith, tt.status)
			}
		})
	}
}
