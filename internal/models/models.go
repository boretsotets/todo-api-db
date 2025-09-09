// Package models реализует структуры, использующиеся в проекте
package models


// Структура PaginatedResponse используется для ответа
// на запрос GET для возвращения значений с учетом нужной страницы
// и размера страницы
type PaginatedResponse struct {
	Data []Task `json:"data"`
	Page int `json:"page"`
	Limit int `json:"limt"`
	Count int `json:"count"`
}

// Структура Task используется при работе с заданиями todo листа
type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedBy *int `json:"createdBy"`
}

// Структура User используется при работе с пользователями todo листа
type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}
