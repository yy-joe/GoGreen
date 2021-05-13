package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"GoGreen/user/data"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Role        string `json:"role" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address"`
	Contact     int    `json:"contact"`
	Date_Joined string `json:"date_joined"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var currentTime = time.Now()
var convertFormat = currentTime.Format("2006-01-02")

var MapUsers map[string]data.User

var MapJSON map[string]UserInfo

func TestData() {
	//users = make(map[string]UserInfo)
	//users["IOT201"] = UserInfo{ID: "0", Username: "nigga", Password: "password", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}

	//mapCourse, _ := json.Marshal(users)
	//fmt.Println("Test Data:", string(mapCourse))
}

// DONE
func LoginAcc(w http.ResponseWriter, r *http.Request) {
	// if user login and username and pw match, add in this isAuthorised function which validate the token in the header
	w.Header().Add("content-type", "application/json")
	db, dbErr := data.ConnectDB()
	if dbErr != nil {
		fmt.Println("error connect to db")
	}
	var login Login

	data.FromJSON(&login, r.Body)
	loginDB, loginErr := data.Login(db, login.Username, login.Password)
	if loginErr != nil {
		fmt.Println("Error retrieving login info from db")
		return
	}

	if login.Username != loginDB.Username {
		fmt.Println("Username do not match")
		w.Write([]byte("401 - Username and/or password do not match."))
		return
	}

	errCompare := bcrypt.CompareHashAndPassword(loginDB.Password, []byte(login.Password))

	if errCompare != nil {
		fmt.Println("Username and/or password do not match", errCompare)
		w.Write([]byte("401 - Username and/or password do not match."))
		return
	}

	data.ToJSON(login, w)
}

// done
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user UserInfo
	db, DBerr := data.ConnectDB()
	if DBerr != nil {
		fmt.Println("Error connect to DB")
	}
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}
	ListUsersFunc()
	if _, ok := MapUsers[user.Username]; ok {
		fmt.Println("Username already taken")
		w.Write([]byte("400 -Username already taken."))
		http.Error(w, "Username already taken", http.StatusForbidden)
		return
	}

	bPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertString := strconv.Itoa(user.Contact)
	user.Date_Joined = convertFormat
	addErr := data.AddUser(db, user.Username, bPassword, user.Name, user.Role, user.Email, user.Address, convertString, user.Date_Joined)
	if addErr != nil {
		fmt.Println("Unable to register new user", addErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error adding new user at database."))
		return
	}
	data.ToJSON(&user, w)
}

func ListUsersFunc() {
	db, err := data.ConnectDB()
	if err != nil {
		fmt.Println("error DB")
	}
	defer db.Close()
	data.GetUsers(db)
	allUsers, err := data.GetUsers(db)
	if err != nil {
		fmt.Println("error getting all users")
		return
	}
	// doing mapping
	MapUsers = make(map[string]data.User)
	MapJSON = make(map[string]UserInfo)
	for _, v := range allUsers {
		MapUsers[string(v.Username)] = v
		MapJSON[v.Username] = UserInfo{v.ID, v.Username, string(v.Password), v.Name, v.Role, v.Email, v.Address, v.Contact, v.Date_Joined}
	}
}

// For admin only
// list all users DONE
func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, err := data.ConnectDB()
	if err != nil {
		fmt.Println("error DB")
	}
	defer db.Close()
	allUsers, err := data.GetUsers(db)
	if err != nil {
		fmt.Println("error getting all users")
		return
	}
	// doing mapping
	MapUsers = make(map[string]data.User)
	for _, v := range allUsers {
		MapUsers[string(v.Username)] = v
	}
	data.ToJSON(allUsers, w)
}

// list userinfo done
func ListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, err := data.ConnectDB()
	if err != nil {
		fmt.Println("error DB")
	}
	defer db.Close()
	params := mux.Vars(r)
	username := params["username"]
	user, err := data.GetUser(db, username)
	data.ToJSON(user, w)
}

// Admin
// Check if username exist or not then add
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, dbErr := data.ConnectDB()
	if dbErr != nil {
		fmt.Println("error connect to db")
	}
	ListUsersFunc()
	var user UserInfo
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}

	if _, ok := MapUsers[user.Username]; ok {
		w.Write([]byte("400 - Username already taken"))
		return
	}

	bPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		w.Write([]byte("400 - Internal server error"))
		return
	}

	convertString := strconv.Itoa(user.Contact)
	user.Date_Joined = convertFormat
	fmt.Println(convertString, bPassword)
	addErr := data.AddUser(db, user.Username, bPassword, user.Name, user.Role, user.Email, user.Address, convertString, user.Date_Joined)
	if addErr != nil {
		w.Write([]byte("400 - Internal server error"))
	}
	data.ToJSON(&user, w)
}

// Done
func EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, dbErr := data.ConnectDB()
	if dbErr != nil {
		fmt.Println("error connect to db")
	}
	var user UserInfo
	params := mux.Vars(r)
	username := params["username"]
	ListUsersFunc()
	data.FromJSON(&user, r.Body)
	convertString := strconv.Itoa(user.Contact)

	bPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	editErr := data.EditUser(db, username, bPassword, user.Name, user.Email, user.Address, convertString)
	if editErr != nil {
		fmt.Println("error DB", editErr)
	}

	if found, ok := MapJSON[params["username"]]; ok {
		MapJSON[params["username"]] = UserInfo{found.ID, found.Username, string(bPassword), user.Name, found.Role, user.Email, user.Address, user.Contact, found.Date_Joined}
		updatedUser := UserInfo{found.ID, found.Username, string(bPassword), user.Name, found.Role, user.Email, user.Address, user.Contact, found.Date_Joined}
		data.ToJSON(updatedUser, w)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - User updated: " +
			params["username"]))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No User found"))
	}
}

// admin done
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, error := data.ConnectDB()
	if error != nil {
		fmt.Println("Error connect DB")
	}
	defer db.Close()
	ListUsersFunc()
	params := mux.Vars(r)
	username := params["username"]
	// validate check
	if _, ok := MapUsers[params["username"]]; ok {
		w.Write([]byte("200 - User deleted"))
		err := data.DeleteUser(db, username)
		if err != nil {
			fmt.Println("error DB", err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No User found"))
	}
}

// done
func GetAdminUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)
	ListUsersFunc()
	if _, ok := MapUsers[params["username"]]; ok {
		w.Write([]byte("200 - User found"))
		json.NewEncoder(w).Encode(MapUsers[params["username"]])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No User found"))
	}
}

// Done
func EditAdminUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, dbErr := data.ConnectDB()
	if dbErr != nil {
		fmt.Println("error connect to db")
	}
	var user UserInfo
	params := mux.Vars(r)
	username := params["username"]
	ListUsersFunc()
	data.FromJSON(&user, r.Body)
	convertString := strconv.Itoa(user.Contact)

	bPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	editErr := data.EditAdminUser(db, username, bPassword, user.Name, user.Role, user.Email, user.Address, convertString)
	if editErr != nil {
		fmt.Println("error DB", editErr)
	}

	if found, ok := MapJSON[params["username"]]; ok {
		MapJSON[params["username"]] = UserInfo{found.ID, found.Username, string(bPassword), user.Name, user.Role, user.Email, user.Address, user.Contact, found.Date_Joined}
		updatedUser := UserInfo{found.ID, found.Username, string(bPassword), user.Name, user.Role, user.Email, user.Address, user.Contact, found.Date_Joined}
		data.ToJSON(updatedUser, w)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - User updated: " +
			params["username"]))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No User found"))
	}
}

func TestAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
}
