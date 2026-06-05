package command

import (
	"errors"
	"testing"
)

type fakeTaskDeleter struct {
	calledWith int
	err        error
}

func (f *fakeTaskDeleter) Delete(id int) error {
	f.calledWith = id
	return f.err
}

func TestDeleteInit(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
		wantID    int
	}{
		{name: "missing id", args: nil, wantError: true},
		{name: "invalid id", args: []string{"abc"}, wantError: true},
		{name: "valid id", args: []string{"3"}, wantID: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCommand(&fakeTaskDeleter{})
			err := cmd.Init(tt.args)

			if tt.wantError && err == nil {
				t.Fatal("Init() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
			if !tt.wantError && cmd.id != tt.wantID {
				t.Fatalf("parsed id = %d, want %d", cmd.id, tt.wantID)
			}
		})
	}
}

func TestDeleteRun(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		storageErr error
		wantError  bool
	}{
		{name: "missing id", id: 0, wantError: true},
		{name: "success", id: 7, wantError: false},
		{name: "storage error propagated", id: 7, storageErr: errors.New("io failure"), wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeTaskDeleter{err: tt.storageErr}
			cmd := NewDeleteCommand(storage)
			cmd.id = tt.id

			err := cmd.Run()
			if tt.wantError && err == nil {
				t.Fatal("Run() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}
			if tt.id > 0 && storage.calledWith != tt.id {
				t.Fatalf("Delete() called with %d, want %d", storage.calledWith, tt.id)
			}
		})
	}
}
