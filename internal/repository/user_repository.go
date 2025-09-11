package repository

import (
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

// UserRepository реализует интерфейс доступа к данным о пользователях.
// Внутри использует подключение к базе данных.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый репозиторий пользователей с переданным
// подключением к базе данных
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// DatabaseInsertUser выполняет SQL-запрос к базе данных для
// сохранения информации о новом пользователе. Возвращает
// id созданного пользователя или ошибку, если запрос не
// может быть обработан
func (r *UserRepository) DatabaseInsertUser(user models.User) (int, error) {
	err := r.db.QueryRow(context.Background(),
		"INSERT INTO users (Email, Name, Password) VALUES ($1, $2, $3) RETURNING Id",
		user.Email, user.Name, user.Password).Scan(&user.Id)
	return user.Id, err
}

// DatabaseRetrieveUser выполняет SQL-запрос к базе данных для
// получения информации о пользователе. Возвращает информацию
// о пользователе или ошибку, если запрос не может быть обработан
func (r *UserRepository) DatabaseRetrieveUser(email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(context.Background(),
		"SELECT Id, Name, Email, Password FROM users WHERE Email = $1",
		email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}
