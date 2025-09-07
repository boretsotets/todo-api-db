package repository

import (
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}


func (r *UserRepository) DatabaseInsertUser(user models.User) (int, error) {
	err := r.db.QueryRow(context.Background(), 
    "INSERT INTO users (Email, Name, Password) VALUES ($1, $2, $3) RETURNING Id", 
	user.Email, user.Name, user.Password).Scan(&user.Id)
	return user.Id, err
}

func (r *UserRepository) DatabaseRetrieveUser(email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(context.Background(), 
    "SELECT Id, Name, Email, Password FROm users WHERE Email = $1", 
	email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}
