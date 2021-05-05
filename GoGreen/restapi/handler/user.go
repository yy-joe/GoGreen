package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goLive/GoGreen/restapi/data"
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
	Date_Joined string `json:"date_joined" validate:"required"`
}

type Login struct {
	Username string
	Password string
}

var currentTime = time.Now()
var convertFormat = currentTime.Format("2006-01-02")

// store all users
var users map[string]UserInfo

func TestData() {
	users = make(map[string]UserInfo)
	users["IOT201"] = UserInfo{ID: "Applied Go Programming", Username: "nigga", Password: "asdasd", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}

	mapCourse, _ := json.Marshal(users)
	fmt.Println("Test Data:", string(mapCourse))
}

func LoginAcc(w http.ResponseWriter, r *http.Request) {
	// if user login and username and pw match, add in this isAuthorised function which validate the token in the header

}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of all courses")
	w.Header().Add("content-type", "application/json")
	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	data.ToJSON(users, w)
	// json.NewEncoder(w).Encode(users)
}

func ListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	users = make(map[string]UserInfo)
	users["test"] = UserInfo{ID: "Nigga Programming", Username: "nigga", Password: "asdasd", Name: "nigga", Role: "admin", Email: "nigga@gmail.com", Contact: 978787, Date_Joined: convertFormat}

	// var user UserInfo
	// if thisUser, ok := users["admin"]; ok {

	// }
	data.ToJSON(users, w)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user UserInfo
	//params := mux.Vars(r)
	//json.NewDecoder(r.Body).Decode(&user)
	data.FromJSON(&user, r.Body)
	// reqBody, err := ioutil.ReadAll(r.Body)
	// if err != nil {

	// }
	// users[params["test"]] = user
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - User added: "))

	//json.NewEncoder(w).Encode(&user)
	data.ToJSON(&user, w)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func getUser() {

}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// var user UserInfo

		// _, found := users[]

		username, password, authOK := r.BasicAuth()
		if authOK == false {
			http.Error(w, "Not authorized", 401)
			return
		}

		if username != "username" || password != "password" {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}
