package models

type PaginatedResponse struct {
	Data []Task `json:"data"`
	Page int `json:"page"`
	Limit int `json:"limt"`
	Count int `json:"count"`
}

type TaskList struct {
	Tasks []Task
}

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedBy *int `json:"createdBy"`
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}
