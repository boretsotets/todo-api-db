package service

import (
	"github.com/boretsotets/todo-api-db/internal/authorization"
	"github.com/boretsotets/todo-api-db/internal/models"
	"github.com/boretsotets/todo-api-db/internal/repository"
	"golang.org/x/crypto/bcrypt"

	"fmt"
)

// UserService содержит бизнес-логику для работы с пользователями.
// Он обращается к репозиторию пользователей для сохранения или получения,
// данных, но сам не зависит от деталей реализации хранилища
type UserService struct {
	repo *repository.UserRepository
}

// NewTaskService создает новый сервис пользователей с внедренным репозиторием
func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

// UserServiceSignIn реализует бизнес-логику авторизации
// пользователя. Принимает авторизационную информацию о
// пользователе. Проверяет существование такого пользователя.
// Возвращает токен авторизации или ошибку, если запрос не может
// быть обработан.
func (s *UserService) UserServiceSignIn(userinfo map[string]string) (string, error) {

	user, err := s.repo.DatabaseRetrieveUser(userinfo["email"])
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userinfo["password"]))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	token, err := authorization.GenerateJWT(user.Id)
	if err != nil {
		return "", fmt.Errorf("token error: %w", err)
	}

	return token, nil
}

// UserServiceSignIn реализует бизнес-логику создания
// нового пользователя. Принимает авторизационную информацию о создаваемом
// пользователе. Хеширует пароль. Передает хеш и информацию
// о пользователе в репозиторий. Возвращает id и авторизационную
// информацию созданного пользователя или ошибку, если запрос
// не может быть обработан.
func (s *UserService) UserServiceSignUp(userinfo map[string]string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(userinfo["password"]), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hash, []byte(userinfo["password"]))
	if err != nil {
		return "", fmt.Errorf("password hashing error: %w", err)
	}
	user := models.User{Name: userinfo["name"], Email: userinfo["email"], Password: string(hash)}

	id, err := s.repo.DatabaseInsertUser(user)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	token, err := authorization.GenerateJWT(id)
	if err != nil {
		return "", fmt.Errorf("token error: %w", err)
	}

	return token, nil
}
