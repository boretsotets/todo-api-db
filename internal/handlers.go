package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
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

type App struct {
	DB *pgxpool.Pool
}

func (a *App) HandlerGet(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	rows, err := a.DB.Query(context.Background(), "SELECT Id, Title, Description FROM tasks")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id int
		var title, description string
		if err := rows.Scan(&id, &title, &description); err != nil {
			log.Fatalf("Rows mapping error: %v", err)
		}	
		tasks = append(tasks, gin.H{"id": id, "title": title, "description": description})
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func (a *App) HandlerPost(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	} 
	
	var newtask JsonTaskPost
	err := json.NewDecoder(c.Request.Body).Decode(&newtask)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
	} else {
		var curr_task Task
		err = a.DB.QueryRow(context.Background(), "INSERT INTO tasks (Title, Description) VALUES ($1, $2) RETURNING (Id, $1, $2)", newtask.Title, newtask.Description).Scan(&curr_task)
		if err != nil {
			log.Fatalf("Rows iteration error: %v", err)
		}
		c.IndentedJSON(http.StatusOK, curr_task)
	}


}

func (a *App) HandlerUpdate(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	} 

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("id conversion error: %v", err)
	}
	var newtask JsonTaskPost
	err = json.NewDecoder(c.Request.Body).Decode(&newtask)

	var curr_task Task
	err = a.DB.QueryRow(context.Background(), "UPDATE tasks SET Title = $1, Description = $2 WHERE Id = $3 RETURNING (Id, Title, Description)", newtask.Title, newtask.Description, id).Scan(&curr_task)
	if err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
	c.IndentedJSON(http.StatusOK, curr_task)

}

func (a *App) HandlerDelete(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if authtoken := c.GetHeader("Authorization"); authtoken == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	} 

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("id conversion error: %v", err)
	}

	a.DB.QueryRow(context.Background(), "DELETE FROM tasks WHERE Id = $1", id)
	c.Status(http.StatusNoContent)
	
}

func (a *App) SignUp(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var userinfo map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&userinfo)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
		log.Fatalf("error decoding JSON: %v", err)
	}
	password := userinfo["password"]

	hash := sha256.New()
	hash.Write([]byte(password))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	password = hashString
	var token string
	err = a.DB.QueryRow(context.Background(), 
    "INSERT INTO users (Email, Name, Password, Token) VALUES ($1, $2, $3, $1) RETURNING Token", userinfo["email"], userinfo["name"], password).Scan(&token)
    if err != nil {
		log.Fatalf("Error querying row: %v", err)
	}
    c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}

func (a *App) SignIn(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var userinfo map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&userinfo)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
		log.Fatalf("error decoding JSON: %v", err)
	}

	password := userinfo["password"]
	hash := sha256.New()
	hash.Write([]byte(password))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	password = hashString

	var token string
	err = a.DB.QueryRow(context.Background(), 
    "SELECT Token FROm users WHERE Email = $1 AND Password = $2", userinfo["email"], password).Scan(&token)
    if err != nil {
		log.Fatalf("Error querying row: %v", err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}
