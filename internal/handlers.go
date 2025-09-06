package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"time"
	"context"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)
//DB_HOST=localhost
//DB_PORT=5432
//DB_USER=postgres
//DB_PASSWORD=secret
//DB_NAME=postgres

// curl -X GET -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcxNTI5NzcsInVzZXJfaWQiOjR9.MaKUr7JTQZNebqMJqID0c6-Xq0ySXYiritC8euROH48" "http://localhost:8080/todos?page=1&limit=10"
// curl -X POST -H "Content-Type application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTcxNTI5NzcsInVzZXJfaWQiOjR9.MaKUr7JTQZNebqMJqID0c6-Xq0ySXYiritC8euROH48" -d '{"title": "title1", "description": "description1"}' http://localhost:8080/todos

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

var jwtKey = []byte("big_secret")

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	return userID, nil
}

func (a *App) HandlerGet(c *gin.Context) {

	authtoken := c.GetHeader("Authorization")
	_, err := ValidateJWT(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		log.Fatalf("page conversion error: %v", err)
	}
	pageSize, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		log.Fatalf("limit conversion error: %v", err)
	}
	offset := (page-1)*pageSize

	rows, err := a.DB.Query(context.Background(), "SELECT Id, Title, Description, CreatedBy FROM tasks")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	index, count := 1, 0
	for rows.Next() {
		if index > offset && count < pageSize {
			var id, createdBy int
			var title, description string
			if err := rows.Scan(&id, &title, &description, &createdBy); err != nil {
				log.Fatalf("Rows mapping error: %v", err)
			}	
			tasks = append(tasks, gin.H{"id": id, "title": title, "description": description})	
			count++
		}
		index++
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": tasks, "page": page, "limit": pageSize, "total": count})
}

func (a *App) HandlerPost(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userID, err := ValidateJWT(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}
	
	var newtask JsonTaskPost
	err = json.NewDecoder(c.Request.Body).Decode(&newtask)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
	} else {
		var curr_task Task
		err = a.DB.QueryRow(context.Background(), "INSERT INTO tasks (Title, Description, CreatedBy) VALUES ($1, $2, $3) RETURNING (Id, $1, $2)", newtask.Title, newtask.Description, userID).Scan(&curr_task)
		if err != nil {
			log.Fatalf("Rows iteration error: %v", err)
		}
		c.IndentedJSON(http.StatusOK, curr_task)
	}


}

func (a *App) HandlerUpdate(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := ValidateJWT(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("id conversion error: %v", err)
	}

	// тут проверка айди клиента и в таске
	var taskOwnerId int
	err = a.DB.QueryRow(context.Background(), "SELECT CreatedBy FROm tasks WHERE Id = $1", id).Scan(&taskOwnerId)
	if err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
	if userId != taskOwnerId {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message" : "Forbidden"})
		return
	}
	

	var newtask map[string]string
	err = json.NewDecoder(c.Request.Body).Decode(&newtask)
	var curr_task Task
	err = a.DB.QueryRow(context.Background(), "UPDATE tasks SET Title = $1, Description = $2 WHERE Id = $3 RETURNING (Id, Title, Description)", newtask["title"], newtask["description"], id).Scan(&curr_task)
	if err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
	c.IndentedJSON(http.StatusOK, curr_task)

}

func (a *App) HandlerDelete(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	authtoken := c.GetHeader("Authorization")
	userId, err := ValidateJWT(authtoken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message" : "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("id conversion error: %v", err)
	}

	// тут проверка айди клиента и в таске
	var taskOwnerId int
	err = a.DB.QueryRow(context.Background(), "SELECT CreatedBy FROm tasks WHERE Id = $1", id).Scan(&taskOwnerId)
	if err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
	if userId != taskOwnerId {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message" : "Forbidden"})
		return
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
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userinfo["password"]), bcrypt.DefaultCost)
	err = bcrypt.CompareHashAndPassword(hash, []byte(userinfo["password"]))
	if err != nil {
		log.Fatalf("error hashing password: %v", err)
	}

	var id int
	err = a.DB.QueryRow(context.Background(), 
    "INSERT INTO users (Email, Name, Password) VALUES ($1, $2, $3) RETURNING Id", 
	userinfo["email"], userinfo["name"], string(hash)).Scan(&id)
    if err != nil {
		log.Fatalf("Error querying row: %v", err)
	}

	token, err := GenerateJWT(id)
	if err != nil {
		log.Fatalf("error generating token: %v", err)
	}
    c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}

func (a *App) SignIn(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var userinfo map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&userinfo)
	if err != nil {
		c.String(http.StatusBadRequest, "Error decoding JSON")
		return
	}

	var id int
	var hashedPassword string
	err = a.DB.QueryRow(context.Background(), 
    "SELECT id, password FROm users WHERE Email = $1", 
	userinfo["email"]).Scan(&id, &hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userinfo["password"]))
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := GenerateJWT(id)
	if err != nil {
		log.Fatalf("error generating token: %v", err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}
