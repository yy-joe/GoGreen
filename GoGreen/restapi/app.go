package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var scanner = bufio.NewScanner(os.Stdin)

const baseURL = "http://localhost:3000/api/v1/"

func addUser() {
	fmt.Println("Please enter the course ID")
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

	jsonData := map[string]string{"id": id, "username": username, "password": password, "name": name, "role": role, "email": email, "address": address, "contact": convertInt}
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(baseURL+"admin/user", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("cannot reach")
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(data))
	response.Body.Close()
}

func main() {
	addUser()
}
