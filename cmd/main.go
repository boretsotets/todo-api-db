package main

import (
	"github.com/gin-gonic/gin"
	"github.com/boretsotets/todo-api-db/internal"
)


func main() {
	router := gin.Default()
	router.GET("/todos", handlers.HandlerGet)
	router.POST("/todos", handlers.HandlerPost)
	router.Run("localhost:8080")

	/*
	tag, err := conn.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT)")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tag)

	tag, err = conn.Exec(context.Background(),
"INSERT INTO tasks (Title, Description) VALUES ('pray isha', 'pray 4 rakaats of isha prayer')")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tag)
	*/


}
