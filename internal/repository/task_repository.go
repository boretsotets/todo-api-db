package repository

import (
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"fmt"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) DatabaseGetTasks(response *models.PaginatedResponse, offset int) error {
	rows, err := r.db.Query(context.Background(), 
	"SELECT Id, Title, Description, CreatedBy FROM tasks LIMIT $1 OFFSET $2",
	response.Limit, offset)
	if err != nil {
		return fmt.Errorf("Database query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var currentTask models.Task
		if err := rows.Scan(&currentTask.Id, &currentTask.Title, &currentTask.Description, &currentTask.CreatedBy); err != nil {
			return fmt.Errorf("Rows mapping error: %w", err)
		}
		response.Data = append(response.Data, currentTask)
		response.Count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("Rows iteration error: %w", err)
	}
	return nil
}

func (r *TaskRepository) DatabaseInsertTask(newTask models.Task) (models.Task, error) {
	var currentTask models.Task
	err := r.db.QueryRow(context.Background(), 
	"INSERT INTO tasks (Title, Description, CreatedBy) VALUES ($1, $2, $3) RETURNING (Id, $1, $2, $3)", 
	newTask.Title, newTask.Description, newTask.Id).Scan(&currentTask)
	return currentTask, err
}

func (r *TaskRepository) DatabaseGetTaskOwner(task models.Task) (int, error) {
	var taskOwnerId int
	err := r.db.QueryRow(context.Background(), 
	"SELECT CreatedBy FROm tasks WHERE Id = $1", task.Id).Scan(&taskOwnerId)
	return taskOwnerId, err
}

func (r *TaskRepository) DatabaseUpdateTask(task models.Task) (models.Task, error) {
	var returnedTask models.Task
	err := r.db.QueryRow(context.Background(), 
	"UPDATE tasks SET Title = $1, Description = $2 WHERE Id = $3 RETURNING (Id, Title, Description)", 
	task.Title, task.Description, task.Id).Scan(&returnedTask)
	return returnedTask, err
}

func (r *TaskRepository) DatabaseRetrieveTaskById(taskId int) (int, error) {
	var taskOwnerId int
	err := r.db.QueryRow(context.Background(), 
	"SELECT CreatedBy FROm tasks WHERE Id = $1", taskId).Scan(&taskOwnerId)
	return taskOwnerId, err
}

func (r *TaskRepository) DatabaseDeleteTask(id int) error {
	var deleteCheckId int
	err := r.db.QueryRow(context.Background(), 
	"DELETE FROM tasks WHERE Id = $1 RETURNING Id", id).Scan(&deleteCheckId)
	return err
}
