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
//DB_HOST=localhost
//DB_PORT=5432
//DB_USER=postgres
//DB_PASSWORD=secret
//DB_NAME=postgres

// curl -X GET -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcyODIzODcsInVzZXJfaWQiOjZ9.4XCU6WWFBt5CG4wL8aa-EOw8q2tNT9ojzCvpOJ9Mre4" "http://localhost:8080/todos?page=6&limit=10"
// curl -X POST -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcyNjY1OTYsInVzZXJfaWQiOjZ9.398bXum_vk8Kk4Vz4bcYS4KnwuzVOZQChmCegpwfBA8" -d '{"title": "title2", "description": "description2"}' http://localhost:8080/todos
// curl -X PUT -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcyNjY1OTYsInVzZXJfaWQiOjZ9.398bXum_vk8Kk4Vz4bcYS4KnwuzVOZQChmCegpwfBA8" -d '{"title": "title12", "description": "description22"}' http://localhost:8080/todos/10
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcyODIzODcsInVzZXJfaWQiOjZ9.4XCU6WWFBt5CG4wL8aa-EOw8q2tNT9ojzCvpOJ9Mre4
// curl -X DELETE -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcyODIzODcsInVzZXJfaWQiOjZ9.4XCU6WWFBt5CG4wL8aa-EOw8q2tNT9ojzCvpOJ9Mre4" http://localhost:8080/todos/55



type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func (h *TaskHandler) HandlerGet(c *gin.Context) {

	authtoken := c.GetHeader("Authorization")
	_, err := authorization.ValidateJWT(authtoken)
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

func (h *TaskHandler) HandlerPost(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var newTask models.Task
	var err error

	authtoken := c.GetHeader("Authorization")
	newTask.Id, err = authorization.ValidateJWT(authtoken)
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

func (h *TaskHandler) HandlerUpdate(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := authorization.ValidateJWT(authtoken)
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
// тут надо выделит файлы, а то методы сломаются

func (h *TaskHandler) HandlerDelete(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := authorization.ValidateJWT(authtoken)
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

