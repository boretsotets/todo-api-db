package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"context"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"strconv"
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

type JsonTaskPost struct {
	Title string`json:"title"`
	Description string`json:"description"`
}

func HandlerGet(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
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
	c.Header("Content-Type", "application/json")

	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		var data RespondUnauthorized
		data.Message = "unauthorized"
		c.IndentedJSON(http.StatusUnauthorized, data)	
	} else {
		var newtask JsonTaskPost
		err := json.NewDecoder(c.Request.Body).Decode(&newtask)
		if err != nil {
			c.String(http.StatusBadRequest, "Error decoding JSON")
		} else {
			conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
			if err != nil {
				log.Fatalf("Database connection failed: %v", err)
			}
			_, err = conn.Exec(context.Background(), 
			"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT)")
			if err != nil {
				log.Fatalf(err)
			}

			var curr_task Task
			err = conn.QueryRow(context.Background(), "INSERT INTO tasks (Title, Description) VALUES ($1, $2) RETURNING (Id, $1, $2)", newtask.Title, newtask.Description).Scan(&curr_task)
			if err != nil {
				log.Fatalf("Rows iteration error: %v", err)
			}
			c.IndentedJSON(http.StatusOK, curr_task)
		}
	}

}

func HandlerUpdate(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		var data RespondUnauthorized
		data.Message = "unauthorized"
		c.Header("Content-Type", "application/json")
		c.IndentedJSON(http.StatusUnauthorized, data)	
	} else {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalf("id conversion error: %v", err)
		}
		var newtask JsonTaskPost
		err = json.NewDecoder(c.Request.Body).Decode(&newtask)

		conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatalf("Database connection failed: %v", err)
		}
		_, err = conn.Exec(context.Background(), 
		"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT)")
		if err != nil {
			log.Fatalf(err)
		}
		var curr_task Task
		err = conn.QueryRow(context.Background(), "UPDATE tasks SET Title = $1, Description = $2 WHERE Id = $3 RETURNING (Id, Title, Description)", newtask.Title, newtask.Description, id).Scan(&curr_task)
		if err != nil {
			log.Fatalf("Rows iteration error: %v", err)
		}
		c.IndentedJSON(http.StatusOK, curr_task)


	}
}

func HandlerDelete(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		var data RespondUnauthorized
		data.Message = "unauthorized"
		c.Header("Content-Type", "application/json")
		c.IndentedJSON(http.StatusUnauthorized, data)	
	} else {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalf("id conversion error: %v", err)
		}
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatalf("Database connection failed: %v", err)
		}
		_, err = conn.Exec(context.Background(), 
		"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Title TEXT, Description TEXT)")
		if err != nil {
			log.Fatalf(err)
		}
		conn.QueryRow(context.Background(), "DELETE FROM tasks WHERE Id = $1", id)
		c.Status(http.StatusNoContent)
	}
}
