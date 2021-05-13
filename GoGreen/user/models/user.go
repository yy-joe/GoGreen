package models

type UserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Role        string `json:"role" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address"`
	Contact     int    `json:"contact"`
	Date_Joined string `json:"date_joined" validate:"required"`
}

type User struct {
	ID          string
	Username    string `validate:"required,alphanum,min=5,max=20"`
	Password    []byte `validate:"required"`
	Name        string `validate:"required,alpha,min=3,max=20"`
	Last        string `validate:"required,alpha,min=3,max=20"`
	Role        string
	Email       string
	Address     string
	Contact     int
	Date_Joined string
}
