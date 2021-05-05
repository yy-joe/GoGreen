package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"
)

// uuid "github.com/satori/go.uuid"
// 	"golang.org/x/crypto/bcrypt"

type loginDetails struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var bPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

// default
var user = loginDetails{
	ID:       1,
	Username: "admin",
	Password: "password",
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

var key = ""
var redisdb *redis.Client
var ctx = context.Background()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key = os.Getenv("SECRET_KEY")

	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisdb = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	// _, err = client.Ping().Result()
	// if err != nil {
	// 	panic(err)
	// }
}

// save meta
func CreateAuth(userid int64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	// rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := redisdb.Set(ctx, td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	// errAccess := client.Set(td.AccessUuid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	// errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	// if errRefresh != nil {
	// 	return errRefresh
	// }
	return nil
}

// Parse, verify and return
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

func home(w http.ResponseWriter, r *http.Request) {
	// validToken, err := CreateToken(1)
	// if err != nil {
	// 	fmt.Println("Invalid token")
	// }

	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}
	r.Header.Set("Token", validToken.AccessToken)

	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Fprintf(w, string(body))
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Endpoint Hit: homePage")
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	var thisUser loginDetails
	err1 := json.Unmarshal(body, &thisUser)
	if err1 != nil {

	}
	// debugging
	fmt.Println(thisUser.Username, thisUser.Password)

	if user.Username != thisUser.Username || user.Password != thisUser.Password {
		// err
		w.Write([]byte("Unauthorised.\n"))
	}

	token, err := GenerateJWT()

	if err != nil {
		w.Write([]byte("Wrong token.\n"))
	}
	fmt.Println("access granted")

	tokens := map[string]string{
		"access_token": token.AccessToken,
	}

	json.NewEncoder(w).Encode(tokens)
}

func GenerateJWT() (*TokenDetails, error) {
	// token details
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	u := uuid.NewV4()
	td.AccessUuid = u.String()
	// td.AccessUuid = uuid.NewV4().String()

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["access_uuid"] = td.AccessUuid
	claims["exp"] = td.AtExpires

	var err error

	td.AccessToken, err = token.SignedString([]byte(key))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return nil, err
	}

	return td, nil
}

func main() {
	// token2, err := GenerateJWT()
	// if err != nil {

	// }
	// fmt.Println(token2)
	router := mux.NewRouter()
	router.Handle("/home", isAuthorized(home)).Methods("GET")
	// router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { isAuthorized(test, home) }).Methods("GET")
	// getR := router.Methods(http.MethodGet).Subrouter()
	// getR.HandleFunc("/home", home).Methods("GET")
	// getR.Use(isAuthorized(test))

	// putR.HandleFunc("/products", ph.Update)
	// putR.Use(ph.MiddlewareValidateProduct)
	router.HandleFunc("/login", Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
