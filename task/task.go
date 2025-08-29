package task
import "time"
type Task struct {
    ID          int
    Title       string
    Description string
    Status      string // "pending", "complete"
    CreatedAt   time.Time
    UpdatedAt   time.Time
}