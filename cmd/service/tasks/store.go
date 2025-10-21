package tasks

import (
	"database/sql"
	"fmt"

	"github.com/Mazin-emad/todo-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}


func (s *Store) CreateTask(task *types.Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (title, completed, user_id) VALUES (?, ?, ?)", task.Title, task.Completed, task.UserID)
	return err
}

func (s *Store) GetTasks(userID int) ([]*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task, err := ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Store) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (s *Store) UpdateTask(task *types.Task) error {
	_, err := s.db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ?", task.Title, task.Completed, task.ID)
	return err
}

func (s *Store) GetTaskByID(id int) (*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	
	task := new(types.Task)
	for rows.Next() {
		task, err = ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	
	if task.ID == 0 {
		return nil, fmt.Errorf("no task found with id: %d", id)
	}
	return task, nil
}

func ScanRowIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)
	err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

