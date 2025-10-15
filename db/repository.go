// Repsitory is just an abstraction over the db so that we don't direcly
//call ORM functions
package db

import "gorm.io/gorm"

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository{
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*User, error) {
    var user User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

type TaskRepository struct {
    db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
    return &TaskRepository{db}
}

func (r *TaskRepository) CreateTask(task *Task) error {
    return r.db.Create(task).Error
}

func (r *TaskRepository) GetTasksByUser(userID uint) ([]Task, error) {
    var tasks []Task
    err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
    return tasks, err
}
func (r *TaskRepository) GetTaskByID(id uint) (*Task, error) {
    var task Task
    err := r.db.First(&task, id).Error
    return &task, err
}

func (r *TaskRepository) UpdateTask(task *Task) error {
    return r.db.Save(task).Error
}

func (r *TaskRepository) DeleteTask(id uint) error {
    return r.db.Delete(&Task{}, id).Error
}

