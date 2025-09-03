package main

import (
	"context"
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/boretsotets/todo-api-db/internal"
)



func main() {
	// подключение к базе данных
	db, err := pgxpool.New(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer db.Close()

	// создание таблицы, если её нет
	_, err = db.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT)")
	if err != nil {
		log.Fatal("error while creating table: ", err)
	}

	_, err = db.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS users (Email text PRIMARY KEY, Name TEXT, Password TEXT, Token TEXT)")
	if err != nil {
		log.Fatal("error while creating table: ", err)
	}

	app := &handlers.App{DB: db}

	router := gin.Default()
	router.POST("/register", app.SignUp)
	router.POST("/login", app.SignIn)
	router.GET("/todos", app.HandlerGet)
	router.POST("/todos", app.HandlerPost)
	router.PUT("/todos/:id", app.HandlerUpdate)
	router.DELETE("todos/:id", app.HandlerDelete)
	router.Run("localhost:8080")

}
