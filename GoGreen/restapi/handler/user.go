package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"GOGREEN/GoGreen/restapi/data"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// uuid "github.com/satori/go.uuid"
//	"golang.org/x/crypto/bcrypt"

// "github.com/gorilla/mux"
type UserInfo struct {
	ID          int    `json:"id"`
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

// store all users
var users map[string]UserInfo

var mapUsers map[string]data.User

var jsonMap map[string]UserInfo

// store cookie session
var mapSessions = map[string]string{}

func getNextID() int {
	// return the last index
	var user UserInfo
	lastIndex := 0
	lastIndex = len(mapUsers)
	fmt.Println("length: ", lastIndex)
	user.ID = lastIndex

	// last index + 1
	return user.ID + 1
}

func TestData() {
	users = make(map[string]UserInfo)
	users["IOT201"] = UserInfo{ID: getNextID(), Username: "nigga", Password: "password", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}

	mapCourse, _ := json.Marshal(users)
	fmt.Println("Test Data:", string(mapCourse))
}

// Search Cookie function
func SetCookie() {

}

// DONE
func LoginAcc(w http.ResponseWriter, r *http.Request) {
	// if user login and username and pw match, add in this isAuthorised function which validate the token in the header
	w.Header().Add("content-type", "application/json")
	var login Login
	_, err := http.Get("https://localhost:3000/api/v1/admin/users")
	if err != nil {
		fmt.Println("error reaching the get all users api", err)
	}

	data.FromJSON(&login, r.Body)

	searchUser, ok := mapUsers[login.Username] // string(bUsername
	if !ok {
		fmt.Println("User not found")
		return
	}
	errCompare := bcrypt.CompareHashAndPassword([]byte(searchUser.Password), []byte(login.Password))

	if errCompare != nil {
		fmt.Println("Username and/or password do not match", errCompare)
		return
	}

	id, _ := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:     "myCookie",
		Value:    id.String(),
		MaxAge:   5,
		HttpOnly: true,
	}

	http.SetCookie(w, myCookie)
	mapSessions[myCookie.Value] = login.Username
	fmt.Println("sessions: ", mapSessions)

	bUsername, _ := bcrypt.GenerateFromPassword([]byte(login.Username), bcrypt.MinCost)
	bPassword, _ := bcrypt.GenerateFromPassword((searchUser.Password), bcrypt.MinCost)
	login.Username = string(bUsername)
	login.Password = string(bPassword)
	data.ToJSON(&login, w)
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user UserInfo
	db, DBerr := data.ConnectDB()
	if DBerr != nil {
		fmt.Println("Error connect to DB")
	}
	_, apiErr := http.Get("https://localhost:3000/api/v1/admin/users")
	if apiErr != nil {
		fmt.Println("error reaching the get all users api", apiErr)
	}
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}
	user.ID = getNextID()

	// need to work on it validate username
	if _, ok := mapUsers[user.Username]; ok {
		fmt.Println("Username already taken")
		http.Error(w, "Username already taken", http.StatusForbidden)
		return
	}

	// create session
	id, _ := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:     "myCookie",
		Value:    id.String(),
		MaxAge:   5,
		HttpOnly: true,
	}
	http.SetCookie(w, myCookie)
	mapSessions[myCookie.Value] = user.Username // cookie become the key

	bPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertString := strconv.Itoa(user.Contact)
	user.Date_Joined = convertFormat
	addErr := data.AddUser(db, strconv.Itoa(user.ID), user.Username, bPassword, user.Name, user.Role, user.Email, user.Address, convertString, user.Date_Joined)
	if addErr != nil {
		fmt.Println("Unable to register new user", addErr)
	}
	fmt.Println("session: ", mapSessions)
	data.ToJSON(&user, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	myCookie, _ := r.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	fmt.Println("Successfully logout")
	http.SetCookie(w, myCookie)
}

// For admin only
// list all users
func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := data.ConnectDB()
	if err != nil {
		fmt.Println("error DB")
	}
	defer db.Close()
	w.Header().Add("content-type", "application/json")
	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	data.GetUsers(db)
	allUsers, err := data.GetUsers(db)
	if err != nil {
		fmt.Println("error getting all users")
		return
	}
	// doing mapping
	mapUsers = make(map[string]data.User)
	for _, v := range allUsers {
		// testMap[string(v.Username)] = UserInfo{v.ID, []byte(v.Username), []byte(v.Password), v.Name, v.Role, v.Email, v.Address, v.Contact, v.Date_Joined}
		// var testMap map[string]data.User
		mapUsers[string(v.Username)] = v
		// testMap[string(v.Username)] = data.User{v.ID}
	}

	// jsonMap = make(map[string]UserInfo)
	// for _, v := range allUsers {
	// 	jsonMap[v.Username] = UserInfo{v.ID, v.Username, []byte(v.Password), v.Name, v.Role, v.Email, v.Address, v.Contact, v.Date_Joined}
	// }
	data.ToJSON(allUsers, w)
	// json.NewEncoder(w).Encode(users)
}

// list userinfo need work on it important
func ListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	myCookie, err := r.Cookie("myCookie")
	if err != nil {
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			HttpOnly: true,
		}

	}
	http.SetCookie(w, myCookie)

	// if the user exists already, get user
	var user UserInfo
	if username, ok := mapSessions[myCookie.Value]; ok {
		user = jsonMap[username]
	}
	data.ToJSON(user, w)
}

// Admin
// Check if username exist or not then add
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user UserInfo
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}
	user.ID = getNextID()
	user.Date_Joined = convertFormat
	data.ToJSON(&user, w)
}

// Admin edit user
func EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var user UserInfo
	params := mux.Vars(r)

	data.FromJSON(&user, r.Body)
	fmt.Print(user.Contact, user.ID)
	if _, ok := users[params["username"]]; ok {
		users[params["username"]] = user
		w.WriteHeader(http.StatusAccepted)
		data.ToJSON(&user, w)
		w.Write([]byte("202 - User updated: " +
			params["username"]))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No User found"))
	}
}

// 	json.NewEncoder(w).Encode(test)
// 	json.NewDecoder(r.Body).Decode(test)
// 	data.FromJSON(&test, r.Body)

// admin done
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db, error := data.ConnectDB()
	if error != nil {
		fmt.Println("Error connect DB")
	}
	// defer db.Close()

	params := mux.Vars(r)
	username := params["username"]
	fmt.Println("username: ", username)
	err := data.DeleteUser(db, username)
	if err != nil {
		fmt.Println("error DB")
		return
	}

	// Local
	// if _, ok := users[params["username"]]; ok {
	// 	fmt.Println("User found")
	// 	delete(users, params["username"])
	// 	w.Write([]byte("202 - User deleted: " +
	// 		params["username"]))
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("404 - No course found"))
	// }
}

// done
func GetAdminUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	params := mux.Vars(r)

	if _, ok := users[params["username"]]; ok {
		fmt.Println("Usesr found")

		json.NewEncoder(w).Encode(users[params["username"]])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found"))
	}
}

func EditAdminUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var user UserInfo
	params := mux.Vars(r)

	data.FromJSON(&user, r.Body)
	if _, ok := users[params["username"]]; ok {
		fmt.Println("user found")
		users[params["username"]] = user
		data.ToJSON(&user, w)
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
	users = make(map[string]UserInfo)
	users["test"] = UserInfo{ID: getNextID(), Username: "nigga", Password: "password", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}
}

func basicAuth(h http.HandlerFunc, Username string, Password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// var user UserInfo

		// _, found := users[]
		db, error := data.ConnectDB()
		if error != nil {
			fmt.Println("Error connect DB")
		}
		defer db.Close()

		// username, password, authOK := r.BasicAuth()
		// if authOK == false {
		// 	http.Error(w, "Not authorized", 401)
		// 	return
		// }
		// login, err := data.Login(db, Username, Password)
		// if err != nil {
		// 	fmt.Println("Error")
		// }
		// if username != login.Username || password != login.Password {
		// 	http.Error(w, "Not authorized", 401)
		// 	return
		// }

		h.ServeHTTP(w, r)
	}
}
