package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var scanner = bufio.NewScanner(os.Stdin)

const baseURL = "https://localhost:3000/api/v1/"

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

var users map[string]UserInfo

func getUser() {

}

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
	users = make(map[string]UserInfo)

	users[username] = UserInfo{id, username, password, name, role, email, address, contact, convertFormat}

	fmt.Println(response.StatusCode)
	fmt.Println(string(data))
	response.Body.Close()
}

func main() {
	addUser(*clientConfig)
	displayAll()
}
