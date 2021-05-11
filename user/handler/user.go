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
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// uuid "github.com/satori/go.uuid"
//	"golang.org/x/crypto/bcrypt"

// "github.com/gorilla/mux"
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

// store all users
var users map[string]UserInfo

var MapUsers map[string]data.User

var jsonMap map[string]UserInfo

// store cookie session
var MapSessions = map[string]string{}

func TestData() {
	users = make(map[string]UserInfo)
	users["IOT201"] = UserInfo{ID: "0", Username: "nigga", Password: "password", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}

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

	id, _ := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:     "myCookie",
		Value:    id.String(),
		MaxAge:   5,
		HttpOnly: true,
	}

	http.SetCookie(w, myCookie)
	MapSessions[myCookie.Value] = login.Username
	fmt.Println("sessions: ", MapSessions)

	bUsername, _ := bcrypt.GenerateFromPassword([]byte(login.Username), bcrypt.MinCost)
	bPassword, _ := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.MinCost)
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
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}
	// sliceUsers, err := data.GetUsers(db)

	// for _, v := range sliceUsers {
	// 	if v.Username == user.Username {
	// 		fmt.Println("Username already taken")
	// 		http.Error(w, "Username already taken", http.StatusForbidden)
	// 		return
	// 	}
	// }
	if _, ok := MapUsers[user.Username]; ok {
		fmt.Println("Username already taken")
		w.Write([]byte("400 -Username already taken."))
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
	MapSessions[myCookie.Value] = user.Username // cookie become the key

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
	fmt.Println("session: ", MapSessions)
	data.ToJSON(&user, w)
}

// func Logout(w http.ResponseWriter, r *http.Request) {
// 	myCookie, _ := r.Cookie("myCookie")
// 	// delete the session
// 	delete(MapSessions, myCookie.Value)
// 	// remove the cookie
// 	myCookie = &http.Cookie{
// 		Name:   "myCookie",
// 		Value:  "",
// 		MaxAge: -1,
// 	}
// 	fmt.Println("Successfully logout")
// 	http.SetCookie(w, myCookie)
// }

// For admin only
// list all users DONE
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
	MapUsers = make(map[string]data.User)
	for _, v := range allUsers {
		// testMap[string(v.Username)] = UserInfo{v.ID, []byte(v.Username), []byte(v.Password), v.Name, v.Role, v.Email, v.Address, v.Contact, v.Date_Joined}
		MapUsers[string(v.Username)] = v
	}

	// jsonMap = make(map[string]UserInfo)
	// for _, v := range allUsers {
	// 	jsonMap[v.Username] = UserInfo{v.ID, v.Username, []byte(v.Password), v.Name, v.Role, v.Email, v.Address, v.Contact, v.Date_Joined}
	// }
	data.ToJSON(allUsers, w)
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
	fmt.Println("cookie", myCookie)
	fmt.Println(MapSessions[myCookie.Value])
	fmt.Println(myCookie.Value)
	// if the user exists already, get user
	var user data.User
	if username, ok := MapSessions[myCookie.Value]; ok {
		fmt.Println("cookie found")
		user = MapUsers[username]
	}
	fmt.Println("user,", user)
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

	var user UserInfo
	data.FromJSON(&user, r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}

	if _, ok := MapUsers[user.Username]; ok {
		fmt.Println("Username already taken")
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
	fmt.Println(convertString, bPassword)
	addErr := data.AddUser(db, user.Username, bPassword, user.Name, user.Role, user.Email, user.Address, convertString, user.Date_Joined)
	if addErr != nil {
		fmt.Println("Unable to register new user", addErr)
		http.Error(w, "Username already taken", http.StatusForbidden)
	}
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
	w.Write([]byte("200 - User deleted"))
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
	users["test"] = UserInfo{ID: "1", Username: "nigga", Password: "password", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}
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
