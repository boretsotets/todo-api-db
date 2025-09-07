package main

import (
	"context"
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/boretsotets/todo-api-db/internal/handlers"
	"github.com/boretsotets/todo-api-db/internal/service"
	"github.com/boretsotets/todo-api-db/internal/repository"

)

// отделение аутентификации
// передача переменных для базы данных и ключа для токена через окружение
// db := setupDB()

func main() {
	// подключение к базе данных
	db, err := pgxpool.New(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer db.Close()

	// создание таблицы, если её нет
	_, err = db.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT, CreatedBy INT)")
	if err != nil {
		log.Fatal("error while creating table: ", err)
	}

	_, err = db.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS users (Id SERIAL PRIMARY KEY, Email TEXT UNIQUE NOT NULL , Name TEXT, Password TEXT)")
	if err != nil {
		log.Fatal("error while creating table: ", err)
	}

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)
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
