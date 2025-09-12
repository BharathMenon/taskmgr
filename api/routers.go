package api

import (
	"strings"
	"log"
	"net/http"
)
type Task struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
//var tasks = []Task{}
func StartServer() {
    http.HandleFunc("/tasks", tasksHandler)
    http.HandleFunc("/tasks/", taskByIDHandler)
    log.Println("API server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        ListTasks(w,r)
        
    case http.MethodPost:
        NewTask(w,r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
	
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
    // route to GetTask, UpdateTask, DeleteTask, MarkComplete
    if strings.HasSuffix(r.URL.Path, "/complete") {
        MarkComplete(w, r)
        return
    }
    switch r.Method {
    case http.MethodGet:
        GetTask(w, r)
    case http.MethodPut:
        UpdateTask(w, r)
    case http.MethodDelete:
        DeleteTask(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}