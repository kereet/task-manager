package database

import (
	"database/sql"
	"fmt"
	"task-manager/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll(ownerID int) ([]models.Task, error) {
	var tasks []models.Task

	query := `
SELECT * 
FROM tasks 
WHERE owner_id = $1
ORDER BY created_at DESC;`

	err := s.db.Select(&tasks, query, ownerID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) GetByID(id, ownerID int) (*models.Task, error) {
	var task models.Task

	query := `
SELECT * 
FROM tasks 
WHERE id = $1 AND owner_id = $2;`

	err := s.db.Get(&task, query, id, ownerID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskStore) Create(input models.CreateTaskInput, ownerID int) (*models.Task, error) {
	var task models.Task
	query := `
INSERT INTO tasks (title, description, completed, owner_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;`

	now := time.Now()
	err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, ownerID, now, now).StructScan(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskStore) Update(id, ownerID int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := s.GetByID(id, ownerID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.UpdatedAt = time.Now()

	query := `
UPDATE tasks
SET title = $1, description = $2, completed = $3, updated_at = $4 
WHERE id = $5 AND owner_id = $6
RETURNING *;`

	var updatedTask models.Task
	err = s.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt, id, ownerID).StructScan(&updatedTask)
	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

func (s *TaskStore) Delete(id, ownerID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND owner_id = $2;`

	result, err := s.db.Exec(query, id, ownerID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}
