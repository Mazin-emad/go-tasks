package types

import "time"

type UserStore interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetAllUsers() ([]*User, error)
}

type TaskStore interface {
	CreateTask(task *Task) error
	GetTasks(userID int) ([]*Task, error)
	DeleteTask(id int) error
	UpdateTask(task *Task) error
	GetTaskByID(id int) (*Task, error)
}

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	UserName     string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUserPayload struct {
	UserName    string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginUserPayload struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateTaskPayload struct {
	Title string `json:"title" validate:"required"`
}

type UpdateTaskPayload struct {
	Title string `json:"title" validate:"required"`
	Completed bool `json:"completed" validate:"required"`
}

type DeleteTaskPayload struct {
	ID int `json:"id" validate:"required"`
}
