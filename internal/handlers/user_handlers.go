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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

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


