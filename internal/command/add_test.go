package command

import (
	"errors"
	"testing"
)

type fakeTaskCreator struct {
	calledWith string
	id         int
	err        error
}

func (f *fakeTaskCreator) Create(description string) (int, error) {
	f.calledWith = description
	if f.err != nil {
		return 0, f.err
	}
	return f.id, nil
}

func TestAddInit(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{name: "missing description", args: nil, wantError: true},
		{name: "has description", args: []string{"buy-milk"}, wantError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeTaskCreator{id: 1}
			cmd := NewAddCommand(storage)

			err := cmd.Init(tt.args)
			if tt.wantError && err == nil {
				t.Fatal("Init() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
		})
	}
}

func TestAddRun(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Add)
		wantError bool
		err       error
		wantDesc  string
	}{
		{
			name: "missing description",
			setup: func(a *Add) {
			},
			wantError: true,
		},
		{
			name: "success",
			setup: func(a *Add) {
				a.description = "ship-v1"
			},
			wantError: false,
			wantDesc:  "ship-v1",
		},
		{
			name: "storage error propagated",
			setup: func(a *Add) {
				a.description = "ship-v1"
			},
			wantError: true,
			err:       errors.New("disk full"),
			wantDesc:  "ship-v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeTaskCreator{id: 11, err: tt.err}
			cmd := NewAddCommand(storage)
			tt.setup(cmd)

			err := cmd.Run()
			if tt.wantError && err == nil {
				t.Fatal("Run() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}
			if tt.wantDesc != "" && storage.calledWith != tt.wantDesc {
				t.Fatalf("Create() called with %q, want %q", storage.calledWith, tt.wantDesc)
			}
		})
	}
}
