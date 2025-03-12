package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserPublic struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
