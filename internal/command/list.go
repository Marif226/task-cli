package command

import (
	"flag"
	"fmt"

	"github.com/Marif226/task-cli/internal/entity"
)

type taskLister interface {
	List(status entity.TaskStatus) ([]entity.Task, error)
}

// List handles the "list" sub-command.
type List struct {
	fs          *flag.FlagSet
	status      string
	taskStorage taskLister
}

// NewListCommand returns a new list command.
func NewListCommand(taskStorage taskLister) *List {
	lc := &List{
		fs:          flag.NewFlagSet("list", flag.ContinueOnError),
		taskStorage: taskStorage,
	}

	return lc
}

// Name returns the sub-command name.
func (l *List) Name() string {
	return l.fs.Name()
}

// Init parses command arguments.
func (l *List) Init(args []string) error {
	err := l.fs.Parse(args)
	if err != nil {
		return err
	}

	if l.fs.NArg() > 0 {
		l.status = l.fs.Arg(0)
	}

	return nil
}

// Run lists the tasks using the parsed status filter if provided.
func (l *List) Run() error {
	tasks, err := l.taskStorage.List(entity.TaskStatus(l.status))
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Println(task.ID, task.Description, string(task.Status))
	}

	return nil
}
