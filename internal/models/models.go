package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}
