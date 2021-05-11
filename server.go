package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"GoGreen/products"
	"GoGreen/user"
	"GoGreen/user/handler"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	handler.TestData()
	// //folder name
	// const folderPath = "templates"
	// tpl = template.Must(template.ParseGlob(folderPath + "/*.html"))

	// //Check templates files exists
	// _, err := ioutil.ReadDir("./templates")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// user.GetAllUsers()
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
	// router.HandleFunc("/api/v1/logout", handler.Logout).Methods("POST")

	// user
	router.HandleFunc("/api/v1/user", handler.ListUser).Methods("GET")
	router.HandleFunc("/api/v1/user/{username}", handler.EditUser).Methods("PUT")

	// admin
	router.HandleFunc("/api/v1/admin/users", handler.ListAllUsers).Methods("GET")
	router.HandleFunc("/api/v1/admin/user", handler.AddUser).Methods("POST")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.GetAdminUser).Methods("GET")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.EditAdminUser).Methods("PUT")
	router.HandleFunc("/api/v1/admin/user/{username}", handler.DeleteUser).Methods("DELETE")

	// GUI template
	// router.HandleFunc("/", user.Index).Methods("GET")
	router.HandleFunc("/login", user.LoginPage).Methods("GET")
	router.HandleFunc("/login", user.UserAuthenticate).Methods("POST")

	router.HandleFunc("/signup", user.SignUpPage).Methods("GET")
	router.HandleFunc("/signup", user.Signup).Methods("POST")
	router.HandleFunc("/logout", user.Logout)

	router.HandleFunc("/user", user.UserPage).Methods("GET")
	router.HandleFunc("/AdminDashboard", user.UserAdminPage).Methods("GET")

	// YYJ & J
	router.HandleFunc("/api/v1/admin/product", products.ProductCRUD).Methods("POST")
	router.HandleFunc("/api/v1/admin/products", products.Allproducts)
	router.HandleFunc("/api/v1/admin/products/active", products.GetActiveProducts)
	router.HandleFunc("/api/v1/admin/products/soldout", products.GetSoldoutProducts)
	router.HandleFunc("/api/v1/admin/products/unlisted", products.GetUnlistedProducts)
	router.HandleFunc("/api/v1/admin/product/quantity-update", products.UpdateProdQty).Methods("PUT")
	router.HandleFunc("/api/v1/admin/product/{productid}", products.ProductCRUD).Methods("GET")
	router.HandleFunc("/api/v1/admin/product/{productid}", products.ProductCRUD).Methods("PUT")
	router.HandleFunc("/api/v1/admin/product/{productid}", products.ProductCRUD).Methods("DELETE")

	router.HandleFunc("/api/v1/admin/brand", products.ServerAddBrand).Methods("POST")
	router.HandleFunc("/api/v1/admin/brands", products.AllBrands)
	router.HandleFunc("/api/v1/admin/brand/{brandid}", products.ServerGetBrand).Methods("GET")
	router.HandleFunc("/api/v1/admin/brand/{brandid}", products.ServerEditBrand).Methods("PUT")
	router.HandleFunc("/api/v1/admin/brand/{brandid}", products.ServerDeleteBrand).Methods("DELETE")

	router.HandleFunc("/api/v1/admin/category", products.ServerAddCategory).Methods("POST")
	router.HandleFunc("/api/v1/admin/categories", products.AllCategories)
	router.HandleFunc("/api/v1/admin/category/{categoryid}", products.ServerGetCategory).Methods("GET")
	router.HandleFunc("/api/v1/admin/category/{categoryid}", products.ServerEditCategory).Methods("PUT")
	router.HandleFunc("/api/v1/admin/category/{categoryid}", products.ServerDeleteCategory).Methods("DELETE")

	router.HandleFunc("/api/v1/admin/enquiry", products.ServerEnquiry).Methods("POST")

	// router.HandleFunc("/api/v1/admin/orders/customer-orders", )
	// router.HandleFunc("/api/v1/admin/orders/product-orders", )

	//handle functions for UI
	//UI URLs for Product Management (Admin)
	router.HandleFunc("/products/all", products.ProdMain)
	router.HandleFunc("/products/{byStatus}", products.ProdByStatus)
	router.HandleFunc("/product/new", products.ProdAdd)
	router.HandleFunc("/product/update/{productid}", products.ProdUpdate)
	router.HandleFunc("/product/delete/{productid}", products.ProdDelete)
	router.HandleFunc("/product/{productid}", products.ProdDetail)

	router.HandleFunc("/enquiry", products.Enquiry)

	//UI URLS for Products/Shop (User)
	router.HandleFunc("/", products.Index)
	router.HandleFunc("/search", products.UserSearch)
	router.HandleFunc("/user/cart", products.Cart)

	//UI URLs for Category Management (Admin)
	router.HandleFunc("/categories/all", products.CatMain)
	router.HandleFunc("/category/new", products.CatAdd)
	router.HandleFunc("/category/{categoryid}", products.CatDetail)
	router.HandleFunc("/category/update/{categoryid}", products.CatUpdate)
	router.HandleFunc("/category/delete/{categoryid}", products.CatDelete)

	//UI URLs for Brand Management (Admin)
	router.HandleFunc("/brands/all", products.BrandMain)
	router.HandleFunc("/brand/new", products.BrandAdd)
	router.HandleFunc("/brand/{brandid}", products.BrandDetail)
	router.HandleFunc("/brand/update/{brandid}", products.BrandUpdate)
	router.HandleFunc("/brand/delete/{brandid}", products.BrandDelete)

	router.HandleFunc("/{productid}", products.Details) //later rename

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
		err := serverConfig.ListenAndServeTLS("cert.pem", "key.pem")
		// err := serverConfig.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

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
	// shutdownContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	serverConfig.Shutdown(shutdownContext)
}
