package command

import (
	"flag"
	"fmt"
)

type taskCreator interface {
	Create(description string) (int, error)
}

// Add handles the "add" sub-command.
type Add struct {
	fs          *flag.FlagSet
	description string
	taskStorage taskCreator
}

// NewAddCommand creates an add command that stores tasks in taskStorage.
func NewAddCommand(taskStorage taskCreator) *Add {
	ac := &Add{
		fs:          flag.NewFlagSet("add", flag.ContinueOnError),
		taskStorage: taskStorage,
	}

	return ac
}

// Name returns the sub-command name.
func (a *Add) Name() string {
	return a.fs.Name()
}

// Init parses command arguments and stores the task description.
func (a *Add) Init(args []string) error {
	err := a.fs.Parse(args)
	if err != nil {
		return err
	}

	if a.fs.NArg() < 1 {
		return fmt.Errorf("description is required")
	}

	a.description = a.fs.Arg(0)

	return nil
}

// Run creates a task using the parsed description.
func (a *Add) Run() error {
	if a.description == "" {
		return fmt.Errorf("description is required")
	}

	id, err := a.taskStorage.Create(a.description)
	if err != nil {
		return err
	}

	fmt.Printf("Task added succesfully (ID: %d)\n", id)

	return nil
}
