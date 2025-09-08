package main

import (
	"context"
	"log"
	
	"github.com/gin-gonic/gin"

	"github.com/boretsotets/todo-api-db/internal/database"
	"github.com/boretsotets/todo-api-db/internal/repository"
	"github.com/boretsotets/todo-api-db/internal/service"
	"github.com/boretsotets/todo-api-db/internal/handlers"
)

// отделение аутентификации
// передача переменных для базы данных и ключа для токена через окружение

func main() {
	// подключение к базе данных
	ctx := context.Background()
	pool, err := database.InitDb(ctx, "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer pool.Close()

	taskRepo := repository.NewTaskRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)
	taskHandler := handlers.NewTaskHandler(taskService)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()

	// user routes
	router.POST("/register", userHandler.HandlerSignUp)
	router.POST("/login", userHandler.HandlerSignIn)

	// task routes
	router.GET("/todos", taskHandler.HandlerGet)
	router.POST("/todos", taskHandler.HandlerPost)
	router.PUT("/todos/:id", taskHandler.HandlerUpdate)
	router.DELETE("todos/:id", taskHandler.HandlerDelete)

	router.Run("localhost:8080")

}
