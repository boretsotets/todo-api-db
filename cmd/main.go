package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"context"
)

//DB_HOST=localhost
//DB_PORT=5432
//DB_USER=postgres
//DB_PASSWORD=secret
//DB_NAME=postgres

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	tag, err := conn.Exec(context.Background(), 
	"CREATE TABLE IF NOT EXISTS tasks (Id SERIAL PRIMARY KEY, Description TEXT, Status TEXT, CreatedAt TIMESTAMP, UpdatedAT TIMESTAMP)")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tag)

	defer conn.Close(context.Background())
}
