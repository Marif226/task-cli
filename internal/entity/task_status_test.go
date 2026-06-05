package entity

import "testing"

func TestTaskStatusIsValid(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{name: "todo is valid", status: TaskStatusToDo, want: true},
		{name: "in-progress is valid", status: TaskStatusInProgress, want: true},
		{name: "done is valid", status: TaskStatusDone, want: true},
		{name: "empty is invalid", status: "", want: false},
		{name: "unknown is invalid", status: "blocked", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTaskStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      TaskStatus
		wantError bool
	}{
		{name: "parse todo", input: "to-do", want: TaskStatusToDo},
		{name: "parse in-progress", input: "in-progress", want: TaskStatusInProgress},
		{name: "parse done", input: "done", want: TaskStatusDone},
		{name: "parse invalid", input: "in review", wantError: true},
		{name: "parse empty", input: "", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTaskStatus(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatalf("ParseTaskStatus() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseTaskStatus() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseTaskStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}
