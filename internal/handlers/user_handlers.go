package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/boretsotets/todo-api-db/internal/service"

	
	"net/http"
	"encoding/json"
	"strings"
)

// curl -X POST -H "Content-Type application/json" -d '{"name": "Donnie Yen", "email": "donnie@yen.com", "password": "donnieyen"}' http://localhost:8080/register
// curl -X POST -H "Content-Type application/json" -d '{"email": "donnie@yen.com", "password": "donnieyen"}' http://localhost:8080/login

// UserHandler реализует HTTP-обработчики для работы с пользователями.
// Он получает зависимость от UserService и использует ее для
// выполнения бизнес-логики.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler создает новый экземпляр UserHandler
// с внедренным сервисом задач
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// HandlerSignUp реализует HTTP-обработчик для работы с POST запросом
// при регистрации нового пользователя. Парсит информацию о пользователе
// и передает управление в сервисный слой. Формирует ответ клиенту в 
// зависимости от результата работы сервиса. В случае успеха возвращает
// сгенерированный токен авторизации
func (h *UserHandler) HandlerSignUp(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var userinfo map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&userinfo)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
		return
	}

	token, err := h.service.UserServiceSignUp(userinfo)
	if err != nil {
		if strings.Contains(err.Error(), "password") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "password hashing error"})
		} else if strings.Contains(err.Error(), "database") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "token error"})
		}
		return
	}

    c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}

// HandlerSignIn реализует HTTP-обработчик для работы с POST запросом
// при логине пользователя. Парсит информацию, необходимую для авторизации и
// передает управление в сервисный слой. Формирует ответ клиенту на основе
// результата работы сервиса. В случае успеха возвращает токен авторизации
func (h *UserHandler) HandlerSignIn(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var userinfo map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&userinfo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error decoding json"})
		return
	}

	token, err := h.service.UserServiceSignIn(userinfo)
	if err != nil {
		if strings.Contains(err.Error(), "database error") {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		} else if strings.Contains(err.Error(), "invalid password") {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid password"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}


