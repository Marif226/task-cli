package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marif226/task-cli/internal/entity"
)

func newTempTaskStorage(t *testing.T, initialContent string) (*Task, string) {
	t.Helper()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "tasks.json")

	if initialContent == "" {
		initialContent = `{"last_id":0,"tasks":[]}`
	}

	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed creating temp storage file: %v", err)
	}

	return NewTaskStorage(filePath), filePath
}

func readTaskStorageFile(t *testing.T, filePath string) *entity.TaskStorage {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed reading storage file: %v", err)
	}

	var storage entity.TaskStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		t.Fatalf("failed unmarshalling storage file: %v", err)
	}

	return &storage
}

func writeTaskStorageFile(t *testing.T, filePath string, storage *entity.TaskStorage) {
	t.Helper()

	data, err := json.Marshal(storage)
	if err != nil {
		t.Fatalf("failed marshalling storage: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatalf("failed writing storage file: %v", err)
	}
}
