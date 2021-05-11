package user

import (
	"GoGreen/user/handler"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"GoGreen/user/data"
)

var scanner = bufio.NewScanner(os.Stdin)

const baseURL = "https://localhost:3000/api/v1/"

const url = "https://localhost:3000/api/v1/login"

var currentTime = time.Now()
var convertFormat = currentTime.Format("2006-01-02")

// transport layer security Configuration
var clientConfig = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type UserInfo struct {
	ID          string
	Username    string
	Password    string
	Name        string
	Role        string
	Email       string
	Address     string
	Contact     int
	Date_Joined string
}

// type LoginInfo struct {
// 	username string `json`
// 	password string `json`
// }

var users map[string]handler.UserInfo
var tpl *template.Template

func init() {
	//folder name
	const folderPath = "user/templates"
	tpl = template.Must(template.ParseGlob(folderPath + "/*.html"))

	//Check templates files exists
	_, err := ioutil.ReadDir("./user/templates")
	if err != nil {
		log.Fatal(err)
	}
	// fetch all users
	// GetAllUsers()
}

func GetAllUsers() {
	// var login Login
	// data.FromJSON(&login, req.Body)
	var users []handler.UserInfo

	response, err := clientConfig.Get(baseURL + "admin/users")
	// data.FromJSON(&user, response.Body)
	if err != nil {
		fmt.Println("", err)
	}
	defer response.Body.Close()
	// fmt.Println("response", response)
	// fmt.Println("username: ", users)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}
	json.Unmarshal(data, &users)
	fmt.Println("username: ", users)
	// users = make(map[string]handler.UserInfo)
	// users[user.Username] = handler.UserInfo{user.ID, user.Username, user.Password, user.Name, user.Role, user.Email, user.Address, user.Contact, user.Date_Joined}

	// for _, v := range users {
	// 	fmt.Println("value: ", v)
	// }
	fmt.Println(response.StatusCode)
	fmt.Println(string(data))
	response.Body.Close()
}

func Index(res http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(res, "index.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

func UserAdminPage(res http.ResponseWriter, req *http.Request) {
	GetAllUsers()
	data := users
	err := tpl.ExecuteTemplate(res, "adminusers.html", data)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

func LoginPage(res http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(res, "login.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

func UserAuthenticate(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	// process form submission
	username := req.FormValue("username")
	password := req.FormValue("password")

	// var login Login
	// data.FromJSON(&login, req.Body)
	newlogin := handler.Login{
		Username: username,
		Password: password,
	}
	jsonValue, _ := json.Marshal(newlogin)
	response, err := clientConfig.Post(baseURL+"login", "application/json", bytes.NewBuffer(jsonValue))
	// fmt.Println("response", response)
	_, err = ioutil.ReadAll(response.Body)

	// json.Unmarshal(reqBody, &newlogin)
	// storedProducts = append(storedProducts, newProduct)
	if err != nil {
		fmt.Println("error from client login", err)
	}
	fmt.Println("Successful login")
	http.Redirect(res, req, "/", http.StatusSeeOther)
	tpl.ExecuteTemplate(res, "login.html", nil)
}

func SignUpPage(res http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(res, "signup.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

func UserPage(res http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(res, "user.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// sign up as a new user
func Signup(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return

	username := req.FormValue("username")
	password := req.FormValue("password")
	name := req.FormValue("name")
	roles := req.FormValue("roles")
	email := req.FormValue("email")
	address := req.FormValue("address")
	contact := req.FormValue("contact")
	convertInt, _ := strconv.Atoi(contact)

	newSignUp := handler.UserInfo{
		Username: username,
		Password: password,
		Name:     name,
		Role:     roles,
		Email:    email,
		Address:  address,
		Contact:  convertInt,
	}

	jsonValue, _ := json.Marshal(newSignUp)

	if username == "" {
		fmt.Println("Username cannot be empty")
		return
	}

	response, err := clientConfig.Post(baseURL+"register", "application/json", bytes.NewBuffer(jsonValue))
	// fmt.Println("response", response)
	resp, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(resp, &newSignUp)
	// json.Unmarshal(reqBody, &newProduct)
	handler.MapUsers = make(map[string]data.User)
	// storedProducts = append(storedProducts, newProduct)
	if err != nil {
		fmt.Println("error from client login", err)
	}
	fmt.Println("Successful register new user")
	http.Redirect(res, req, "/", http.StatusSeeOther)
	tpl.ExecuteTemplate(res, "signup.html", newSignUp)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	fmt.Println("LOGOUT function")
	// if !alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }
	fmt.Println("cookie", handler.MapSessions)
	myCookie, _ := req.Cookie("myCookie")
	http.Redirect(res, req, "/login", http.StatusSeeOther)
	// delete the session
	delete(handler.MapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	fmt.Println("logout successfully")
	fmt.Println(handler.MapSessions)

}

// Add new user
// func AddUser(res http.ResponseWriter, req *http.Request) {
// 	myUser := getUser(res, req)
// 	//var myUser user

// 	if myUser.Role == "admin" {
// 		tpl.ExecuteTemplate(res, "addUser.html", myUser)
// 		// process form submission
// 		if req.Method == http.MethodPost {
// 			// get form values
// 			addUsername := req.FormValue("addUsername")
// 			addPassword := req.FormValue("addPassword")

// 			addfirstname := req.FormValue("addfirstname")
// 			addLastname := req.FormValue("addLastname")
// 			addRoles := req.FormValue("addRoles")
// 			var nameAlpha = true
// 			for _, char := range addfirstname {
// 				if unicode.IsLetter(char) == false {
// 					nameAlpha = false
// 				}
// 			}

// 			var nameLength = false
// 			if len(addfirstname) < 20 {
// 				nameLength = true
// 			}

// 			if !nameAlpha || !nameLength {
// 				http.Error(res, "Error input, please check your input", http.StatusForbidden)
// 				return
// 			}

// 			if addUsername != "" {
// 				// check if username exist/ taken
// 				if _, ok := mapUsers[addUsername]; ok {
// 					http.Error(res, "Username already taken", http.StatusForbidden)
// 					return
// 				}

// 				bPassword, err := bcrypt.GenerateFromPassword([]byte(addPassword), bcrypt.MinCost)
// 				if err != nil {
// 					http.Error(res, "Internal server error", http.StatusInternalServerError)
// 					return
// 				}

// 				myUser := user{addUsername, bPassword, addfirstname, addLastname, addRoles}
// 				err3 := validateUser.Struct(myUser)
// 				if err3 != nil {
// 					fmt.Println("error")
// 					http.Error(res, "error input, please check", http.StatusForbidden)
// 					return
// 				}
// 				mapUsers[addUsername] = myUser
// 			}
// 			return
// 		}
// 	} else {
// 		tpl.ExecuteTemplate(res, "noAccess.html", myUser)
// 	}
// }

// func getUser(res http.ResponseWriter, req *http.Request) user {
// 	// get current session cookie
// 	myCookie, err := req.Cookie("myCookie")
// 	if err != nil {
// 		id, _ := uuid.NewV4()
// 		myCookie = &http.Cookie{
// 			Name:     "myCookie",
// 			Value:    id.String(),
// 			HttpOnly: true,
// 		}

// 	}
// 	http.SetCookie(res, myCookie)

// 	// if the user exists already, get user
// 	var myUser user
// 	if username, ok := mapSessions[myCookie.Value]; ok {
// 		myUser = mapUsers[username]
// 	}
// 	return myUser
// }

func displayAll() {
	for _, v := range users {

		fmt.Println("display all", v)
	}
}

func editUser() {
	fmt.Println("Please enter the ID")
	var id string
	fmt.Scanln(&id)

	fmt.Println("Please enter the username")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Please enter the password")
	var password string
	fmt.Scanln(&password)

	fmt.Println("Please enter the name")
	var name string
	fmt.Scanln(&name)

	fmt.Println("Please enter the role")
	var role string
	fmt.Scanln(&role)

	fmt.Println("Please enter the email")
	var email string
	fmt.Scanln(&email)

	fmt.Println("Please enter the address")
	var address string
	fmt.Scanln(&address)

	fmt.Println("Please enter the int")
	var contact int
	fmt.Scanln(&contact)

	convertInt := strconv.Itoa(contact)

	jsonData := map[string]string{"id": id, "username": username, "password": password, "name": name, "role": role, "email": email, "address": address, "contact": convertInt, "date_joined": convertFormat}
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.NewRequest(baseURL+"admin/user", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("cannot reach")
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	var user UserInfo
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}

	// fmt.Println(response.StatusCode)
	json.Unmarshal(data, &user)
	fmt.Println(user)
	// users = append(users, user)
	fmt.Println(string(data))
	response.Body.Close()
}

func addUser(client http.Client) {
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan() // use `for scanner.Scan()` to keep reading
	// id := scanner.Text()

	fmt.Println("Please enter the ID")
	var id string
	fmt.Scanln(&id)

	fmt.Println("Please enter the username")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Please enter the password")
	var password string
	fmt.Scanln(&password)

	fmt.Println("Please enter the name")
	var name string
	fmt.Scanln(&name)

	fmt.Println("Please enter the role")
	var role string
	fmt.Scanln(&role)

	fmt.Println("Please enter the email")
	var email string
	fmt.Scanln(&email)

	fmt.Println("Please enter the address")
	var address string
	fmt.Scanln(&address)

	fmt.Println("Please enter the int")
	var contact int
	fmt.Scanln(&contact)

	convertInt := strconv.Itoa(contact)

	jsonData := map[string]string{"id": id, "username": username, "password": password, "name": name, "role": role, "email": email, "address": address, "contact": convertInt, "date_joined": convertFormat}
	jsonValue, _ := json.Marshal(jsonData)

	response, err := client.Post(baseURL+"admin/user", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("cannot reach")
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}
	// users = make(map[string]handler.UserInfo)

	// users[username] = handler.UserInfo{id, username, password, name, role, email, address, contact, convertFormat}

	fmt.Println(response.StatusCode)
	fmt.Println(string(data))
	response.Body.Close()
}

// func main() {
// 	addUser(*clientConfig)
// 	displayAll()
// }
