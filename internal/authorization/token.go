// Package authorization реализует функции для генерации и валидации
// JWT токена при авторизации пользователей
package authorization

import (
	"github.com/golang-jwt/jwt/v5"

	"os"
	"time"
)

// Секретный пароль, использующийся для
// подписи JWT токенов
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT генерирует новый токен авторизации для пользователя.
// Включает идентификатор пользователя и время жизни токена - 1 час. В
// Используется при авторизации после входа в систему. Возвращает токен или ошибку
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT проверяет, является ли tokenString валидным токеном
// авторизации пользователя. Если токен действителен, возвращает
// идентификатор пользователя userID. Если токен недействителен или
// истек, возвращает ошибку
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
