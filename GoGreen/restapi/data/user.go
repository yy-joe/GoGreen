package data

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// connect to database

// Data model
type User struct {
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

// Users defines a slice of User
type Users []*User

var currentTime = time.Now()
var convertFormat = currentTime.Format("01-02-2006")

// default admin acc
var userList = []*User{
	&User{
		ID:          "1",
		Username:    "admin",
		Password:    "password",
		Name:        "admin",
		Role:        "admin",
		Email:       "admin.gmail.com",
		Address:     "",
		Contact:     0,
		Date_Joined: convertFormat,
	},
}

// Return a list of all users
func GetUsers() {

}

// Return
func GetUser() {

}

//Authentication
func Login() {

}
