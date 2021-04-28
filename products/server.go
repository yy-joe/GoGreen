package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	Trace *log.Logger //Prints execution status to stdout, for debugging purposes
)

func openDB() (db *sql.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			Trace.Println("Panic trapped: ", err)
			return
		}
	}()

	//Use mysql as driverName and the default mysql db as data source name
	dsn := "root:QQ2kepiting@tcp(127.0.0.1:3306)/GoGreen"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Trace.Fatalln(err.Error())
	}
	return
}

func allproducts(w http.ResponseWriter, r *http.Request) {

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	products, err := getProducts(db)
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting products from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(products)
}

func product(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	if r.Method == "GET" {
		if product, err := getProduct(db, params["productid"]); err != nil { //check if productid exists in the database
			Trace.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No product found."))
		} else {
			//w.WriteHeader(http.StatusOK)
			//w.Write([]byte("200 - Found requested product."))
			//json.NewEncoder(w).Encode(product)
			fmt.Println(product)
		}
	} else if r.Method == "DELETE" {
		productID, _ := strconv.Atoi(params["productid"])
		err := deleteProducts(db, productID)
		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error deleting product from database."))
			return
		}
		// w.WriteHeader(http.StatusAccepted)
		// w.Write([]byte("202 - product deleted: " + params["productid"]))
		fmt.Println("Product deleted:", productID)

		// } else if r.Header.Get("Content-type") == "application/json" {

		//POST is for creating new product
	} else if r.Method == "POST" {
		//read the string sent to the service
		var newProduct Product
		// reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply product information in JSON format"))
			return
		} else {
			//convert JSON to object
			// json.Unmarshal(reqBody, &newProduct)

			// if newProduct.Name == "" {
			// 	w.WriteHeader(http.StatusUnprocessableEntity)
			// 	w.Write([]byte("422 - Please supply product information in JSON format"))
			// 	return
			// }

			newProduct = Product{
				ID:         0,
				Name:       "Prod123",
				Image:      "",
				Details:    "Test product",
				DateAdded:  "2021-04-28",
				Price:      10.00,
				Quantity:   50,
				CategoryID: 2,
				BrandID:    2,
			}
			//check if product exists; add only if product does not exist
			err := addProducts(db, newProduct.Name, newProduct.Image, newProduct.Details, newProduct.DateAdded, newProduct.Price, newProduct.Quantity, newProduct.CategoryID, newProduct.BrandID)
			if err != nil {
				Trace.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 -Error updating product at database."))
				return
			}
			// w.WriteHeader(http.StatusCreated)
			// w.Write([]byte("201 - product added: " + params["productid"]))
			fmt.Println("Product successfully added.")
		}
	} else if r.Method == "PUT" { //PUT is for creating or updating existing product
		var newProduct Product
		// reqBody, err := ioutil.ReadAll(r.Body)

		newProduct = Product{
			ID:         0,
			Name:       "Prod123",
			Image:      "",
			Details:    "Test product",
			DateAdded:  "2021-04-28",
			Price:      8.00,
			Quantity:   150,
			CategoryID: 2,
			BrandID:    2,
		}

		//product already exists, update product
		productID, _ := strconv.Atoi(params["productid"])
		err := editProducts(db, newProduct.Name, newProduct.Image, newProduct.Details, newProduct.DateAdded, newProduct.Price, newProduct.Quantity, newProduct.CategoryID, newProduct.BrandID, productID)

		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error updating product at database."))
			return
		}
		// w.WriteHeader(http.StatusAccepted)
		// w.Write([]byte("202 - product updated: " + params["productid"]))

	}
}

func main() {

	router := mux.NewRouter()
	// router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/products", allproducts)
	router.HandleFunc("/api/v1/product/{productid}", product).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/api/v1/product", product).Methods("POST")

	fmt.Println("Listening at port 5000")
	//log.Fatal(http.ListenAndServe(":5000", router))

	err := http.ListenAndServe(":5000", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
