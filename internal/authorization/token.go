// Package authorization реализует функции для генерации и валидации
// JWT токена при авторизации пользователей
package authorization

import (
	"github.com/golang-jwt/jwt/v5"

	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Generate JWT генерирует JWT токен для пользователя с идентификатором
// userId, действующий в течение одного часа. Возвращает токен и ошибку
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT проверяет, является ли tokenString валидным токеном
// авторизации пользователя. Возвращает userID пользователя, для которого
// был сгенерирован токен, и ошибку.
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
