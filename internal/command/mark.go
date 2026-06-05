package command

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/Marif226/task-cli/internal/entity"
)

type taskStatusUpdater interface {
	UpdateStatus(id int, status entity.TaskStatus) (*entity.Task, error)
}

// Mark handles the "mark" sub-command.
type Mark struct {
	fs          *flag.FlagSet
	id          int
	status      entity.TaskStatus
	taskStorage taskStatusUpdater
}

// NewMarkCommand returns a new mark command.
func NewMarkCommand(taskStorage taskStatusUpdater) *Mark {
	mc := &Mark{
		fs:          flag.NewFlagSet("mark", flag.ContinueOnError),
		taskStorage: taskStorage,
	}

	return mc
}

// Name returns the sub-command name.
func (m *Mark) Name() string {
	return m.fs.Name()
}

// Init parses command arguments and stores the task description.
func (m *Mark) Init(args []string) error {
	err := m.fs.Parse(args)
	if err != nil {
		return err
	}

	if m.fs.NArg() < 2 {
		return fmt.Errorf("not enough arguments: task ID and status are required")
	}

	id, err := strconv.Atoi(m.fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	m.id = id

	status, err := entity.ParseTaskStatus(m.fs.Arg(1))
	if err != nil {
		return err
	}

	m.status = status

	return nil
}

// Run marks a task using the parsed task ID and status.
func (m *Mark) Run() error {
	if m.id <= 0 {
		return fmt.Errorf("task ID must be greater than 0")
	}

	_, err := m.taskStorage.UpdateStatus(m.id, m.status)
	if err != nil {
		return err
	}

	fmt.Printf("Task marked successfully (ID: %d, Status: %s)\n", m.id, string(m.status))

	return nil
}
