package main

import (
	"github.com/gin-gonic/gin"
	"github.com/boretsotets/todo-api-db/internal"
)


func main() {
	router := gin.Default()
	router.GET("/todos", handlers.HandlerGet)
	router.POST("/todos", handlers.HandlerPost)
	router.PUT("/todos/:id", handlers.HandlerUpdate)
	router.DELETE("todos/:id", handlers.HandlerDelete)
	router.Run("localhost:8080")

}
