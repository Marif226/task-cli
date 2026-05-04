package storage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Marif226/task-cli/internal/entity"
)

// Task stores tasks in a JSON file.
type Task struct {
	filePath string
}

// NewTaskStorage creates file-backed task storage.
func NewTaskStorage(filePath string) *Task {
	return &Task{filePath: filePath}
}

// Create adds a new task and returns its ID.
func (t *Task) Create(description string) (int, error) {
	taskStorage, err := t.loadStorage()
	if err != nil {
		return 0, err
	}

	task := &entity.Task{
		ID:          taskStorage.LastID + 1,
		Status:      entity.TaskStatusToDo,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	taskStorage.Tasks = append(taskStorage.Tasks, *task)
	taskStorage.LastID++

	err = t.atomicSave(taskStorage)
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (t *Task) loadStorage() (*entity.TaskStorage, error) {
	data, err := os.ReadFile(t.filePath)
	if err != nil {
		return nil, err
	}

	var taskStorage entity.TaskStorage

	if len(data) != 0 {
		err = json.Unmarshal(data, &taskStorage)
		if err != nil {
			return nil, err
		}
	}

	return &taskStorage, nil
}

func (t *Task) atomicSave(taskStorage *entity.TaskStorage) error {
	data, err := json.Marshal(taskStorage)
	if err != nil {
		return err
	}

	tmpFile := t.filePath + ".tmp"
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		return err
	}

	return os.Rename(tmpFile, t.filePath)
}
