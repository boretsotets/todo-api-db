package service

import (
	"github.com/boretsotets/todo-api-db/internal/authorization"
	"github.com/boretsotets/todo-api-db/internal/repository"
	"github.com/boretsotets/todo-api-db/internal/models"

	"fmt"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(r *repository.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) ServiceGet(response *models.PaginatedResponse) (error) {
	offset := (response.Page-1)*response.Limit
	err := s.repo.DatabaseGetTasks(response, offset)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *TaskService) ServicePost(newTask models.Task) (models.Task, error){
	currentTask, err := s.repo.DatabaseInsertTask(newTask)
	if err != nil {
		return currentTask, fmt.Errorf("error inserting task: %w", err)
	}
	return currentTask, nil
}

func (s *TaskService) ServiceUpdate(task models.Task, userId int) (models.Task, error) {
	var newTask models.Task
	// тут проверка айди клиента и в таске
	taskOwnerId, err := s.repo.DatabaseGetTaskOwner(task)
	if err != nil {
		return newTask, fmt.Errorf("database retrieval error: %w", err)
	}
	if userId != taskOwnerId {
		return newTask, fmt.Errorf("owner rights error: %w", err)
	}

	newTask, err = s.repo.DatabaseUpdateTask(task)
	if err != nil {
		return newTask, fmt.Errorf("task update error: %w", err)
	}
	return newTask, nil
}

func (s *TaskService) ServiceDelete(userId, id int) (error) {
	var taskOwnerId int
	taskOwnerId, err := s.repo.DatabaseRetrieveTaskById(id)
	if err != nil {
		return fmt.Errorf("Task retrieving error: %w", err)
	}
	if userId != taskOwnerId {
		return fmt.Errorf("User don't have rights to delete this task: %w", err)
		// c.IndentedJSON(http.StatusForbidden, gin.H{"message" : "Forbidden"})
	}
	err = s.repo.DatabaseDeleteTask(id)
	if err != nil {
		return fmt.Errorf("Database error: %w", err)
	}
	return nil
}

func ServiceAuth(token string) (int, error) {
	userId, err := authorization.ValidateJWT(token)
	return userId, err
}