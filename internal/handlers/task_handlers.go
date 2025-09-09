// Package handlers отвечает за транспортный слой приложения.
// Здесь находятся HTTP-обработчики, которые принимают запросы,
// передают управление в сервисный слой и формируют ответы для клиента.
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/boretsotets/todo-api-db/internal/service"
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/boretsotets/todo-api-db/internal/authorization"


	"net/http"
	"encoding/json"
	"strconv"
	"strings"
)

// curl -X GET -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTczNjYxNTYsInVzZXJfaWQiOjZ9.YXtBeLRVQCKTQXwPVC0nLvIzN_rIowYFuhX-cUeI8Jc" "http://localhost:8080/todos?page=6&limit=10"
// curl -X POST -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTc0NDgzODIsInVzZXJfaWQiOjd9.XM2LwiReRdwMzAGUorHC0_hDAOnU2PnKq8-F9N65RnA" -d '{"title": "title3", "description": "description3"}' http://localhost:8080/todos
// curl -X PUT -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTc0NDY5NzIsInVzZXJfaWQiOjJ9.GbVaoCciaXWT70VY67PmnOLzFzD9v52OnndWTGsFhOU" -d '{"title": "title22", "description": "description22"}' http://localhost:8080/todos/1
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTczNjYxNTYsInVzZXJfaWQiOjZ9.YXtBeLRVQCKTQXwPVC0nLvIzN_rIowYFuhX-cUeI8Jc
// curl -X DELETE -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTc0NDY5NzIsInVzZXJfaWQiOjJ9.GbVaoCciaXWT70VY67PmnOLzFzD9v52OnndWTGsFhOU" http://localhost:8080/todos/55


// TaskHandler реализует HTTP-обработчики для работы с todo-задачами.
// Он получает зависимость от TaskService и использует ее для
// выполнения бизнес-логики.
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler создает новый экземпляр TaskHandler
// с внедренным сервисом задач
func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// HabdlerGet реализует HTTP-обработчик для работы с запросом GET.
// Проверяет токен авторизации пользователя, парсит query параметры
// page и limit. После этого передает управление в сервисный слой
// и формирует ответ клиенту в зависимости от результата работы сервиса.
func (h *TaskHandler) HandlerGet(c *gin.Context) {

	authtoken := c.GetHeader("Authorization")
	_, err := service.ServiceAuth(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	var response models.PaginatedResponse
	response.Page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		c.String(http.StatusBadRequest, "error converting page to integer")
		return
	}
	response.Limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.String(http.StatusBadRequest, "error converting limit to integer")
		return
	}
	err = h.service.ServiceGet(&response)
	if err != nil {
		if strings.Contains(err.Error(), "query") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database query error"})
		} else if strings.Contains(err.Error(), "mapping") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "rows mapping error"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "rows iteration error"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

// HandlerPost реализует HTTP-обработчик для работы с запросом POST
// при создании новых todo-тасков. Проверяет авторизацию пользователя,
// парсит информацию о новой задаче и передает управление в сервисный слой.
// Формирует ответ клиенту в зависимости от результата работы сервиса.
func (h *TaskHandler) HandlerPost(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var newTask models.Task
	var err error

	authtoken := c.GetHeader("Authorization")
	newTask.Id, err = service.ServiceAuth(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}
	
	err = json.NewDecoder(c.Request.Body).Decode(&newTask)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
		return
	}
	currentTask, err := h.service.ServicePost(newTask)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error inserting task"})
		return
	}

	c.IndentedJSON(http.StatusOK, currentTask)
	

}

// HandlerUpdate реализует HTTP-обработчик для работы с запросом PUT
// при изменении todo-тасков. Проверяет авторизацию пользователя, 
// парсит id задачи, которую надо изменить, и новые поля для задачи.
// Передает управление в сервисный слой и формирует ответ клиенту 
// на основе результата работы сервиса.
func (h *TaskHandler) HandlerUpdate(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := service.ServiceAuth(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	var oldTask models.Task
	oldTask.Id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error converting id to integer"})
		return
	}

	err = json.NewDecoder(c.Request.Body).Decode(&oldTask)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error converting request body to json"})
		return
	}

	newTask, err := h.service.ServiceUpdate(oldTask, userId)
	if err != nil {
		if strings.Contains(err.Error(), "database") {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error retrieving this task from database"})
		} else if strings.Contains(err.Error(), "owner") {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "you don't have rights to change this task"})
		} else if strings.Contains(err.Error(), "task") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error updating task"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": newTask.Id, "title": newTask.Title, "description": newTask.Description})

}

// HandlerDelete реализует HTTP-обработчик для запроса DELETE
// при удалении todo-таска. Проверяет авторизацию пользователя, 
// парсит id задачи, которую надо удалить и передает управление
// в сервисный слой. Формирует ответ клиенту на основе работы сервиса.
func (h *TaskHandler) HandlerDelete(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := service.ServiceAuth(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message" : "id convertion error"})
		return
	}

	// тут проверка айди клиента и в таске
	err = h.service.ServiceDelete(userId, id)
	if err != nil {
		if strings.Contains(err.Error(), "retrieving") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving task"})
		} else if strings.Contains(err.Error(), "rights") {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "user have no rights to delete this task"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error deleting task"})
		}
		return
	}
	
	c.Status(http.StatusNoContent)
	
}

