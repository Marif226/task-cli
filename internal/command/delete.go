package command

import (
	"flag"
	"fmt"
	"strconv"
)

type taskDeleter interface {
	Delete(id int) error
}

// Delete handles the "delete" sub-command.
type Delete struct {
	fs          *flag.FlagSet
	id          int
	taskStorage taskDeleter
}

// NewDeleteCommand creates a delete command that deletes tasks from storage.
func NewDeleteCommand(taskStorage taskDeleter) *Delete {
	rc := &Delete{
		fs:          flag.NewFlagSet("delete", flag.ContinueOnError),
		taskStorage: taskStorage,
	}

	return rc
}

// Name returns the sub-command name.
func (r *Delete) Name() string {
	return r.fs.Name()
}

// Init parses command arguments and stores the task description.
func (r *Delete) Init(args []string) error {
	err := r.fs.Parse(args)
	if err != nil {
		return err
	}

	if r.fs.NArg() < 1 {
		return fmt.Errorf("task ID is required")
	}

	id, err := strconv.Atoi(r.fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	r.id = id

	return nil
}

// Run deletes a task using the parsed description.
func (r *Delete) Run() error {
	if r.id == 0 {
		return fmt.Errorf("task ID is required")
	}

	err := r.taskStorage.Delete(r.id)
	if err != nil {
		return err
	}

	return nil
}
