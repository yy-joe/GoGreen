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

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var scanner = bufio.NewScanner(os.Stdin)

const baseURL = "https://localhost:3000/api/v1/"

const adminBaseURL = "https://localhost:3000/api/v1/admin/user/"

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
	Password    []byte
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

var mapUsers map[string]data.User
var tpl *template.Template

// store cookie session
var MapSessions = map[string]string{}

func init() {
	//folder name
	const folderPath = "user/templates"
	tpl = template.Must(template.ParseGlob(folderPath + "/*.html"))

	//Check templates files exists
	_, err := ioutil.ReadDir("./user/templates")
	if err != nil {
		log.Fatal("file or folder doesn't exists", err)
	}
}

// check if user logged in done
func AlreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := MapSessions[myCookie.Value]
	// search for
	_, ok := mapUsers[username]
	return ok
}

// done
func GetAllUsers() {
	var users []data.User
	response, err := clientConfig.Get(baseURL + "admin/users")
	if err != nil {
		fmt.Println("", err)
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}
	json.Unmarshal(respBody, &users)
	mapUsers = make(map[string]data.User)
	for _, v := range users {
		mapUsers[v.Username] = v
	}
	response.Body.Close()
}

// index page done
func Index(res http.ResponseWriter, req *http.Request) {
	user := GetUser(res, req)
	err := tpl.ExecuteTemplate(res, "index.html", user)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// done
func GetUser(res http.ResponseWriter, req *http.Request) data.User {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		fmt.Println("error getting cookie", err)
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			HttpOnly: true,
		}

	}
	http.SetCookie(res, myCookie)
	// if the user exists already, get user
	var myUser data.User
	if username, ok := MapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}
	return myUser
}

// done
func UserAdminPage(res http.ResponseWriter, req *http.Request) {
	GetAllUsers()
	data := mapUsers
	err := tpl.ExecuteTemplate(res, "adminusers.html", data)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// done
func LoginPage(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(res, "login.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// done
func UserAuthenticate(res http.ResponseWriter, req *http.Request) {
	GetAllUsers()
	// process form submission
	username := req.FormValue("username")
	password := req.FormValue("password")

	newlogin := handler.Login{
		Username: username,
		Password: password,
	}

	searchUser, ok := mapUsers[username]
	if !ok {
		http.Error(res, "Username do not match", http.StatusUnauthorized)
		return
	}
	err := bcrypt.CompareHashAndPassword(searchUser.Password, []byte(password))
	if err != nil {
		http.Error(res, "Username and/or password do not match", http.StatusForbidden)
		return
	}

	jsonValue, _ := json.Marshal(newlogin)
	response, err := clientConfig.Post(baseURL+"login", "application/json", bytes.NewBuffer(jsonValue))
	_, err = ioutil.ReadAll(response.Body)

	id, _ := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:     "myCookie",
		Value:    id.String(),
		MaxAge:   5,
		HttpOnly: true,
	}

	http.SetCookie(res, myCookie)
	MapSessions[myCookie.Value] = username
	// json.Unmarshal(reqBody, &newlogin)
	// storedProducts = append(storedProducts, newProduct)
	if err != nil {
		fmt.Println("error from client login", err)
	}
	fmt.Println("Successful login")
	http.Redirect(res, req, "/", http.StatusSeeOther)
	tpl.ExecuteTemplate(res, "login.html", nil)
}

// done
func SignUpPage(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(res, "signup.html", nil)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// done
func UserPage(res http.ResponseWriter, req *http.Request) {
	user := GetUser(res, req)
	// var user UserInfo
	response, err := clientConfig.Get(baseURL + "api/v1/user/" + user.Username)
	if err != nil {
		fmt.Println("", err)
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}
	json.Unmarshal(respBody, &user)

	response.Body.Close()
	// fmt.Println(user.Username)
	err = tpl.ExecuteTemplate(res, "user.html", user)
	if err != nil {
		http.Error(res, "Server Time out", http.StatusRequestTimeout)
	}
}

// Done sign up as a new user
func Signup(res http.ResponseWriter, req *http.Request) {
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
	resp, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(resp, &newSignUp)

	if err != nil {
		fmt.Println("error from client login", err)
	}
	fmt.Println("Successful register new user")
	http.Redirect(res, req, "/", http.StatusSeeOther)
	tpl.ExecuteTemplate(res, "signup.html", newSignUp)
}

//done
func Logout(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("logout function")
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(MapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	fmt.Println("logout successsful")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// done
func EditUser(res http.ResponseWriter, req *http.Request) {
	user := GetUser(res, req)

	if req.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(res, "editUser.html", user)
		if err != nil {
			fmt.Println(err)
		}
	}

	if req.Method == http.MethodPost {
		newpassword := req.FormValue("newpassword")
		confirmpassword := req.FormValue("confirmpassword")
		name := req.FormValue("name")
		roles := req.FormValue("roles")
		email := req.FormValue("email")
		address := req.FormValue("address")
		contact := req.FormValue("contact")
		convertInt, _ := strconv.Atoi(contact)

		if newpassword != confirmpassword {
			fmt.Println("password")
			return
		}

		updateUser := handler.UserInfo{
			Username:    user.Username,
			Password:    confirmpassword,
			Name:        name,
			Role:        roles,
			Email:       email,
			Address:     address,
			Contact:     convertInt,
			Date_Joined: user.Date_Joined,
		}
		jsonValue, _ := json.Marshal(updateUser)
		request, err := http.NewRequest(http.MethodPut, baseURL+"user/"+user.Username, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
		request.Header.Set("Content-Type", "application/json")
		response, errConfig := clientConfig.Do(request)
		if errConfig != nil {
			log.Fatalln(err)
		}

		_, reserr := ioutil.ReadAll(response.Body)
		if reserr != nil {
			fmt.Println("reserr", reserr)
		}
		// json.Unmarshal(reqBody, &updateUser)
		fmt.Println("after", updateUser)
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}

// done
func EditAdminUser(res http.ResponseWriter, req *http.Request) {
	user := GetUser(res, req)

	if req.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(res, "editUser.html", user)
		if err != nil {
			fmt.Println(err)
		}
	}

	if req.Method == http.MethodPost {
		newpassword := req.FormValue("newpassword")
		confirmpassword := req.FormValue("confirmpassword")
		name := req.FormValue("name")
		roles := req.FormValue("roles")
		email := req.FormValue("email")
		address := req.FormValue("address")
		contact := req.FormValue("contact")
		convertInt, _ := strconv.Atoi(contact)

		if newpassword != confirmpassword {
			fmt.Println("password")
			return
		}

		updateUser := handler.UserInfo{
			Username:    user.Username,
			Password:    confirmpassword,
			Name:        name,
			Role:        roles,
			Email:       email,
			Address:     address,
			Contact:     convertInt,
			Date_Joined: user.Date_Joined,
		}
		jsonValue, _ := json.Marshal(updateUser)
		request, err := http.NewRequest(http.MethodPut, baseURL+"user/"+user.Username, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
		request.Header.Set("Content-Type", "application/json")
		response, errConfig := clientConfig.Do(request)
		if errConfig != nil {
			log.Fatalln(err)
		}

		_, reserr := ioutil.ReadAll(response.Body)
		if reserr != nil {
			fmt.Println("reserr", reserr)
		}
		// json.Unmarshal(reqBody, &updateUser)
		fmt.Println("after", updateUser)
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}

// admin Done delete
func DeleteUser(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		deleteuser := req.FormValue("delete")

		req, err := http.NewRequest(http.MethodDelete, adminBaseURL+deleteuser, nil)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		response, err := clientConfig.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if response.StatusCode != 200 {
			return
		}

		//direct user back to the main products page
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	name := req.FormValue("name")
	roles := req.FormValue("roles")
	email := req.FormValue("email")
	address := req.FormValue("address")
	contact := req.FormValue("contact")
	convertInt, _ := strconv.Atoi(contact)

	newUser := handler.UserInfo{
		Username: username,
		Password: password,
		Name:     name,
		Role:     roles,
		Email:    email,
		Address:  address,
		Contact:  convertInt,
	}
	fmt.Println("Add new user:", newUser)

	jsonValue, _ := json.Marshal(newUser)

	if username == "" {
		fmt.Println("Username cannot be empty")
		return
	}
	response, err := clientConfig.Post(baseURL+"admin/user", "application/json", bytes.NewBuffer(jsonValue))
	resp, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(resp, &newUser)

	if err != nil {
		fmt.Println("error from client login", err)
	}
	fmt.Println("Successful add new user")
	http.Redirect(res, req, "/", http.StatusSeeOther)
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
