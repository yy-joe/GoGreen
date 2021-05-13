package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// connect to database

// Data model
type User struct {
	ID          string
	Username    string
	Password    []byte
	Name        string
	Role        string
	Email       string
	Address     string
	Contact     int
	Date_Joined string
}

// Users defines a slice of User
type Users []*User

var currentTime = time.Now()
var convertFormat = currentTime.Format("01-02-2006")

var Trace *log.Logger

// default admin acc
var adminTest = []*User{
	&User{
		ID:          "1",
		Username:    "admin",
		Password:    []byte("password"),
		Name:        "admin",
		Role:        "admin",
		Email:       "admin.gmail.com",
		Address:     "",
		Contact:     0,
		Date_Joined: convertFormat,
	},
}

func ConnectDB() (db *sql.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			Trace.Println("Panic trapped: ", err)
			return
		}
	}()
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/GoGreen")
	// handle error
	if err != nil {
		fmt.Println("error open mysql")
		// panic(err.Error())
	}
	fmt.Println("Database Connected")
	return
}

// Return a list of all users
func GetUsers(db *sql.DB) ([]User, error) {
	results, err := db.Query("SELECT * FROM GoGreen.Users")
	var users []User
	if err != nil {
		fmt.Println("error get results", err)
		return nil, err
		// panic(err.Error())
	}
	defer results.Close()
	for results.Next() { // map this type to the record in the table
		var user User
		err = results.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Role, &user.Email, &user.Address, &user.Contact, &user.Date_Joined)
		if err != nil {
			fmt.Println("error scanning", err)
			return nil, err
			// panic(err.Error())
		}
		users = append(users, user)
	}
	return users, nil
}

// Done
func GetUser(db *sql.DB, username string) (User, error) {
	var user User

	err := db.QueryRow("SELECT * FROM Users WHERE Username=?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Role, &user.Email, &user.Address, &user.Contact, &user.Date_Joined)

	// error connecting to mysql DB
	if err != nil {
		fmt.Println("Error", err)
		panic(err.Error())
	}
	return user, nil
}

// Add User
func AddUser(db *sql.DB, Username string, Password []byte, Name string, Role string, Email string, Address string, Contact string, DateJoined string) error {
	query := fmt.Sprintf("INSERT INTO Users (Username, Password, Name, Role, Email, Address, Contact, Date_Joined) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", Username, Password, Name, Role, Email, Address, Contact, DateJoined)
	_, err := db.Exec(query)
	if err != nil {
		Trace.Println(err)
		fmt.Println("Error adding new user in mysql", err)
	}
	return err
}

func EditUser(db *sql.DB, Username string, Password []byte, Name string, Email string, Address string, Contact string) error {
	query := fmt.Sprintf("UPDATE Users SET Password='%s', Name='%s', Email='%s', Address='%s', Contact='%s' WHERE Username='%s'", Password, Name, Email, Address, Contact, Username)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error", err)
		panic(err.Error())
	}
	return err
}

func EditAdminUser(db *sql.DB, Username string, Password []byte, Name string, Role string, Email string, Address string, Contact string) error {
	query := fmt.Sprintf("UPDATE Users SET Password='%s', Name='%s', Role='%s', Email='%s', Address='%s', Contact='%s' WHERE Username='%s'", Password, Name, Role, Email, Address, Contact, Username)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error", err)
		panic(err.Error())
	}
	return err
}

func DeleteUser(db *sql.DB, Username string) error {
	query := fmt.Sprintf("DELETE FROM Users WHERE Username='%s'", Username)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Unable to delete", err)
		panic(err.Error())
	}
	return err
}

//Authentication
func Login(db *sql.DB, Username string, Password string) (User, error) {
	var user User
	err := db.QueryRow("SELECT * FROM Users WHERE Username=?", Username).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Role, &user.Email, &user.Address, &user.Contact, &user.Date_Joined)
	if err != nil {
		fmt.Println("Cannot find username", err)
		// log.Fatalln(err)
		// panic(err.Error())
	}
	return user, nil
}
