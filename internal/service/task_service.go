// Package service реализует бизнес-логику приложения.
// Сервисы инкапсулируют сценарии использования,
// работают с репозиторием и пакетом авторизации.
package service

import (
	"github.com/boretsotets/todo-api-db/internal/authorization"
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/boretsotets/todo-api-db/internal/repository"

	"fmt"
)

// TaskService содержит бизнес-логику для работы с todo-задачами.
// Он обращается к репозиторию задач для сохранения, получения,
// изменения или удаления данных, но сам не зависит от деталей
// реализации хранилища
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService создает новый сервис todo-задач с внедренным репозиторием
func NewTaskService(r *repository.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

// ServiceGet обрабатывает бизнес-логику получения списка задач.
// На основе информации о размере страницы и номера нужной страницы,
// считает необходимый отступ, а затем запрашивает данные у репозитория.
// Возвращает список задач и информацию о пагинации или ошибку,
// если запрос не может быть выполнен
func (s *TaskService) ServiceGet(response *models.PaginatedResponse) error {
	offset := (response.Page - 1) * response.Limit
	err := s.repo.DatabaseGetTasks(response, offset)
	if err != nil {
		return err
	}

	return nil
}

// ServicePost обрабатывает бизнес-логику создания новой todo-задачи.
// Принимает информацию о новой задаче и передает ее репозиторию.
// Возвращает поля созданной задачи или ошибку.
func (s *TaskService) ServicePost(newTask models.Task) (models.Task, error) {
	currentTask, err := s.repo.DatabaseInsertTask(newTask)
	if err != nil {
		return currentTask, fmt.Errorf("error inserting task: %w", err)
	}
	return currentTask, nil
}

// ServiceUpdate обрабатывает бизнес-логику изменения todo-задачи.
// Принимает id задачи, которую надо изменить, и изменяемые поля.
// Затем обновляет поля задачи через репозиторий.
// Возвращает поля измененной задачи или ошибку, если запрос не
// может быть обработан.
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

// ServiceDelete обрабатывает бизнес-логику удаления todo-задачи.
// Принимает id клиента, отправившего запрос, и id удаляемой задачи.
// Проверяет права клиента на удаление этой задачи а затем удаляет ее
// через репозиторий. В случае, если запрос не может быть обработан,
// возвращает ошибку
func (s *TaskService) ServiceDelete(userId, id int) error {
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

// ServiceAuth обрабатывает бизнес-логику авторизации пользователя.
// Принимает токен авторизации и проверяет существование пользователя
// с таким действующим токеном через пакет авторизации.
// Возвращает id пользователя или ошибку, если такого пользователя нет.
func ServiceAuth(token string) (int, error) {
	userId, err := authorization.ValidateJWT(token)
	return userId, err
}
