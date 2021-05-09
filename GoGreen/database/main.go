package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// default
var user = User{
	ID:       1,
	Username: "admin",
	Password: "password",
}

var key = ""

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key = os.Getenv("SECRET_KEY")
}

// var mySigningKey = []byte("c78afaf-97da-4816-bbee-9ad239abb296")

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(key), nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Println("Not authorized")
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

// func loginFunc(x string, y string) map[string]interface {

// }

func testtest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Endpoint Hit: homePage")
}

func home(w http.ResponseWriter, r *http.Request) {
	// validToken, err := CreateToken(1)
	// if err != nil {
	// 	fmt.Println("Invalid token")
	// }

	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}
	r.Header.Set("Token", validToken)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(body))
	fmt.Fprintf(w, validToken)
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Endpoint Hit: homePage")
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	var thisUser User
	err1 := json.Unmarshal(body, &thisUser)
	if err1 != nil {

	}

	// debugging
	fmt.Println(thisUser.Username, thisUser.Password)
	// login := loginFunc(user.Username, user.Password)
	if user.Username != thisUser.Username || user.Password != thisUser.Password {
		// err
		w.Write([]byte("Unauthorised.\n"))
	}

	token, err := GenerateJWT()
	if err != nil {
		w.Write([]byte("Wrong token.\n"))
	}
	fmt.Println("access granted")
	json.NewEncoder(w).Encode(token)
}

func GenerateJWT() (string, error) {
	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["user_id"] = userid
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func main() {
	fmt.Println(key)
	fmt.Println(mySigningKey)
	// test, err := CreateToken()
	// if err != nil {

	// }
	// fmt.Println(test)

	token2, err := GenerateJWT()
	if err != nil {

	}
	fmt.Println(token2)
	router := mux.NewRouter()
	router.Handle("/home", isAuthorized(home)).Methods("GET")
	router.HandleFunc("/test", testtest).Methods("GET")
	// router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { isAuthorized(test, home) }).Methods("GET")
	// getR := router.Methods(http.MethodGet).Subrouter()
	// getR.HandleFunc("/home", home).Methods("GET")
	// getR.Use(isAuthorized(test))

	// putR.HandleFunc("/products", ph.Update)
	// putR.Use(ph.MiddlewareValidateProduct)
	router.HandleFunc("/login", Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
