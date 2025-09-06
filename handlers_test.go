package handlers

import "testing"

func TestHandlerPost(t *testing.T) {
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
	
		app := &handlers.App{DB: db}
	
}
