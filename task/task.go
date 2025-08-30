package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	//"io/ioutil"
	"os"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string // "pending", "complete"
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const (
	StatusPending = "pending"
	StatusDone    = "done"
)

func (t *Task) MarkComplete() {
	t.Status = StatusDone
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// fallback to UTC if IST can't be loaded
		t.UpdatedAt = time.Now().UTC()
		return
	}
	t.UpdatedAt = time.Now().In(ist)
}

func  TasksFilePath() string {
	if p := os.Getenv("TASK_FILE"); p != "" {
		return p
	}
	return "tasks.json"
}

func  loadTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // no file => empty list
		}
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(path string, tasks []Task) error {
	// Marshal pretty to help debugging
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		dir = "."
	}
	tmpFile, err := os.CreateTemp(dir, "tasks_*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		// best-effort
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}
	// Atomic replace
	return os.Rename(tmpName, path)
}

func nextID(tasks []Task) int {
	var max int = 0
	for _, t := range tasks {
		if (t.ID) > max {
			max = (t.ID)
		}
	}
	return max + 1
}

func  AddTask(path, title, desc string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("title is required")
	}
	tasks, err := loadTasks(path)
	if err != nil {
		return Task{}, err
	}
	ist, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        ist = time.UTC
    }
    now := time.Now().In(ist)
	task := Task{
		ID:          nextID(tasks),
		Title:       title,
		Description: desc,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, task)
	if err := saveTasks(path, tasks); err != nil {
		return Task{}, err
	}
	return task, nil
}

func ListTasks(path string) ([]Task, error) {
	tasks, err := loadTasks(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks, nil
}

func findTaskIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func (t *Task) Update(title, desc, status *string) {
	now := time.Now().UTC()
	if title != nil {
		t.Title = *title
	}
	if desc != nil {
		t.Description = *desc
	}
	if status != nil {
		t.Status = *status
	}
	t.UpdatedAt = now
}

func UpdateTask(path string, id int, title, desc, status *string) (Task, error) {
	if id <= 0 {
		return Task{}, errors.New("id must be greater than 0")
	}
	tasks, err := loadTasks(path)
	if err != nil {
		return Task{}, err
	}
	idx := findTaskIndex(tasks, id)
	if idx < 0 {
		return Task{}, fmt.Errorf("task %d not found", id)
	}
	tasks[idx].Update(title, desc, status)
	if err := saveTasks(path, tasks); err != nil {
		return Task{}, err
	}
	return tasks[idx], nil
}

func DeleteTask(path string, id int) error {
	if id <= 0 {
		return errors.New("id must be > 0")
	}
	tasks, err := loadTasks(path)
	if err != nil {
		return err
	}
	idx := findTaskIndex(tasks, id)
	if idx < 0 {
		return fmt.Errorf("task %d not found", id)
	}
	tasks = append(tasks[:idx], tasks[idx+1:]...) //... or ellipsis unpacks tasks[idx+1:] to individual elements
	return saveTasks(path, tasks)
}

func CompleteTask(path string, id int) (Task, error) {
	tasks, err := loadTasks(path)
	if err != nil {
		return Task{}, err
	}
	idx := findTaskIndex(tasks, id)
	if idx < 0 {
		return Task{}, fmt.Errorf("task %d not found", id)
	}
	tasks[idx].MarkComplete()
	if err := saveTasks(path, tasks); err != nil {
		return Task{}, err
	}
	return tasks[idx], nil
}

func PrintTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks.")
		return
	}
	for _, t := range tasks {
    status := t.Status
    fmt.Printf("ID: %d | %s | %s\n", t.ID, t.Title, status)
    if t.Description != "" {
        fmt.Printf("   %s\n", t.Description)
    }
    ist, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        ist = time.UTC
    }
    fmt.Printf(
        "   Created: %s Updated: %s\n",
        t.CreatedAt.In(ist).Format(time.RFC3339),
        t.UpdatedAt.In(ist).Format(time.RFC3339),
    )
}
}