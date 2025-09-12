package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/BharathMenon/taskmgr/task"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := task.ListTasks(task.TasksFilePath())
	if err != nil {
		http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") //Whenver I am sending json data back to client I use application/json
	json.NewEncoder(w).Encode(tasks)
}

// GetTask - GET /tasks/{id}
func GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	t, err := task.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// UpdateTaskHandler - PUT /tasks/{id}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"desc"`
		Status      *string `json:"status"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	t, err := task.UpdateTask(task.TasksFilePath(), id, req.Title, req.Description, req.Status)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// DeleteTaskHandler - DELETE /tasks/{id}
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	err = task.DeleteTask(task.TasksFilePath(), id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
// MarkComplete - PUT /tasks/{id}/complete
func MarkComplete(w http.ResponseWriter,r *http.Request){
if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr:= strings.TrimPrefix(r.URL.Path,"/tasks/")
	idStr = strings.TrimSuffix(idStr,"/complete")
	id,err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }
	path:=task.TasksFilePath()
	t,err := task.CompleteTask(path,id)
	if err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
	t.Status = "done"
	 w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(t)

}

// func ListTasks(w http.ResponseWriter, r *http.Request) {
// 	tasks, err := task.ListTasks(task.TasksFilePath())
// 	if err != nil {
// 		http.Error(w, "Failed to List Tasks", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(tasks)
// }

func NewTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"desc"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	t, err := task.AddTask(task.TasksFilePath(), req.Title, req.Description)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// func GetTask(w http.ResponseWriter, r *http.Request){

// }