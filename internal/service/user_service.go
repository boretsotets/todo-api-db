package service

import (
	"github.com/boretsotets/todo-api-db/internal/authorization"
	"github.com/boretsotets/todo-api-db/internal/repository"
	"github.com/boretsotets/todo-api-db/internal/models"
	"golang.org/x/crypto/bcrypt"

	"fmt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

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
