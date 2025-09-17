package api
import (
	"net/http"
    "strconv"

    task "github.com/BharathMenon/taskmgr/task"
    gin "github.com/gin-gonic/gin"
)
type NewTaskRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"desc"`
}

func ListTasks(c *gin.Context) {
    tasks, err := task.ListTasks(task.TasksFilePath())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
        return
    }
    c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    t, err := task.GetTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, t)
}

func UpdateTask(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    if c.GetHeader("Content-Type") != "application/json" {
        c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be application/json"})
        return
    }
    var req struct {
        Title       *string `json:"title"`
        Description *string `json:"desc"`
        Status      *string `json:"status"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }
    t, err := task.UpdateTask(task.TasksFilePath(), id, req.Title, req.Description, req.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }
    c.JSON(http.StatusOK, t)
}

func DeleteTask(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    if err := task.DeleteTask(task.TasksFilePath(), id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.Status(http.StatusNoContent)
}

func MarkComplete(c *gin.Context) {
    // With Gin, declare route as /tasks/:id/complete so "id" is available as a param.
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    t, err := task.CompleteTask(task.TasksFilePath(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    // Optionally override status if needed
    t.Status = "done"
    c.JSON(http.StatusOK, t)
}

func NewTask(c *gin.Context) {
    // Create a binding struct using the defined NewTaskRequest
    var req NewTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }
    // Title is required by binding:"required" tag.
    t, err := task.AddTask(task.TasksFilePath(), req.Title, req.Description)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }
    c.JSON(http.StatusCreated, t)
}