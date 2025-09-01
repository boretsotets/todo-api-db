package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"context"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
)
//DB_HOST=localhost
//DB_PORT=5432
//DB_USER=postgres
//DB_PASSWORD=secret
//DB_NAME=postgres

type TaskList struct {
	Tasks []Task
}

type Task struct {
	Id int
	Title string
	Description string
}

func HandlerGet(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	rows, err := conn.Query(context.Background(), "SELECT Id, Title, Description FROM tasks")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	data := TaskList{}
	for rows.Next() {
		var curr_task Task
		
		if err := rows.Scan(&curr_task.Id, &curr_task.Title, &curr_task.Description); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		data.Tasks = append(data.Tasks, curr_task)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
	c.IndentedJSON(http.StatusOK, data)

	defer conn.Close(context.Background())
}

type RespondUnauthorized struct {
	Message string`json:"message"`
}

func HandlerPost(c *gin.Context) {
	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		var data RespondUnauthorized
		data.Message = "unauthorized"
		c.Header("Content-Type", "application/json")
		c.IndentedJSON(http.StatusUnauthorized, data)	
	}

}
