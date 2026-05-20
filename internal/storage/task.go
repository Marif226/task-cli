package storage

import (
	"encoding/json"
	"fmt"
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

// Delete deletes a task by its ID.
func (t *Task) Delete(id int) error {
	taskStorage, err := t.loadStorage()
	if err != nil {
		return err
	}

	// Find the task by ID and delete it from the slice.
	for i, task := range taskStorage.Tasks {
		if task.ID == id {
			taskStorage.Tasks = append(taskStorage.Tasks[:i], taskStorage.Tasks[i+1:]...)
			return t.atomicSave(taskStorage)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

// loadStorage reads the task storage from the JSON file and returns it as a TaskStorage struct.
// If the file does not exist, it returns an empty TaskStorage struct.
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

// atomicSave writes the task storage to a temporary file and then renames it to the original file path.
// This ensures that the file is updated atomically, preventing data corruption in case of a crash during the write operation.
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
