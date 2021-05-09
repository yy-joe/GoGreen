package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"goLive/GoGreen/restapi/handler"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	handler.TestData()
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	log.Println("Hello nigga")
	rw.Header().Set("Content-Type", "application/json")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Fprintf(rw, "hello %s", d)
}

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portNo := os.Getenv("PORT")

	l := log.New(os.Stdout, "course-api", log.LstdFlags)

	router := mux.NewRouter()
	// auth
	router.HandleFunc("/api/v1/login", handler.LoginAcc).Methods("POST")
	router.HandleFunc("/api/v1/register", handler.Register).Methods("POST")
	router.HandleFunc("/api/v1/logout", handler.Logout).Methods("POST")

	// user
	router.HandleFunc("/api/v1/user", handler.ListUser).Methods("GET")
	router.HandleFunc("/api/v1/user/{username}", handler.EditUser).Methods("PUT")

	// admin
	router.HandleFunc("/api/v1/admin/users", handler.ListAllUsers).Methods("GET")
	router.HandleFunc("/api/v1/admin/user", handler.AddUser).Methods("POST")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.GetAdminUser).Methods("GET")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.EditAdminUser).Methods("PUT")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.DeleteUser).Methods("DELETE")

	// getRouter := router.Methods("GET").Subrouter()
	// getRouter.HandleFunc("/api/v1/admin/user/{username}", handler.GetAdminUser)

	// putRouter := router.Methods("PUT").Subrouter()
	// putRouter.HandleFunc("/api/v1/admin/user/{username}", handler.EditUser)

	// server config
	serverConfig := http.Server{
		Addr:         ":" + portNo,      // configure the bind address
		Handler:      router,            // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  5 * time.Second,   // max time to read request from the client
		ReadTimeout:  10 * time.Second,  // max time to write response to the client
		WriteTimeout: 120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port", portNo)
		// err := serverConfig.ListenAndServeTLS("cert.pem", "key.pem")
		err := serverConfig.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// db, err := handler.ConnectDB()
	// data.GetUsers(db)
	// serverConfig.ListenAndServe()
	// Signal notification channel
	// signalChannel := make(chan os.Signal)
	signalChannel := make(chan os.Signal, 1)

	// broadcast msg into signal channel
	signal.Notify(signalChannel, os.Interrupt)
	// signal.Notify(signalChannel, os.Kill)
	signal.Notify(signalChannel, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-signalChannel
	fmt.Println("Recieve terminate, graceful shutdow", sig)
	log.Println("Recieve terminate, graceful shutdow", sig)

	// gracefully shutdown and close all the workers in 30sec
	shutdownContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	serverConfig.Shutdown(shutdownContext)

}
