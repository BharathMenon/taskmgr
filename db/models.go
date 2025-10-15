package db
import(
	"time"
)
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Tasks     []Task    `gorm:"constraint:OnDelete:CASCADE"`
    CreatedAt time.Time
}

type Task struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Description string
    Status      string    `gorm:"default:'pending'"`
    UserID      uint      `gorm:"not null"`
    User        User      `gorm:"foreignKey:UserID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

