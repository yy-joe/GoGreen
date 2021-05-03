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
	dsn := "root:password@tcp(127.0.0.1:3306)/GoGreen"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Trace.Fatalln(err.Error())
	}
	return
}

func activeProducts(w http.ResponseWriter, r *http.Request) {
	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	products, err := getProductsByStatus(db, "active")
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting products from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(products)
}

func soldoutProducts(w http.ResponseWriter, r *http.Request) {
	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	products, err := getProductsByStatus(db, "soldout")
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting products from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(products)
}

func unlistedProducts(w http.ResponseWriter, r *http.Request) {
	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	products, err := getProductsByStatus(db, "unlisted")
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting products from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(products)
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
				ID:           0,
				Name:         "Prod123",
				Image:        "",
				DescShort:    "Test product",
				DescLong:     "Long Test product",
				DateCreated:  "",
				DateModified: "",
				Price:        10.00,
				Quantity:     50,
				Condition:    "New",
				CategoryID:   2,
				BrandID:      3,
				Status:       "Live",
			}
			//check if product exists; add only if product does not exist
			err := addProducts(db, newProduct.Name, newProduct.Image, newProduct.DescShort, newProduct.DescLong, newProduct.Price, newProduct.Quantity, newProduct.Condition, newProduct.CategoryID, newProduct.BrandID, newProduct.Status)
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
		// reqBody, err := ioutil.ReadAll(r.Body)

		newProduct := Product{
			ID:           0,
			Name:         "Prod123",
			Image:        "",
			DescShort:    "Test product",
			DescLong:     "Long Test product",
			DateCreated:  "",
			DateModified: "",
			Price:        20.00,
			Quantity:     30,
			Condition:    "New",
			CategoryID:   2,
			BrandID:      3,
			Status:       "Sold Out",
		}

		//product already exists, update product
		productID, _ := strconv.Atoi(params["productid"])
		err := editProducts(db, newProduct.Name, newProduct.Image, newProduct.DescShort, newProduct.DescLong, newProduct.Price, newProduct.Quantity, newProduct.Condition, newProduct.CategoryID, newProduct.BrandID, newProduct.Status, productID)

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

func allBrands(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	brands, err := getBrands(db)
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting products from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(brands)
}

func serverGetBrand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the brand database."))
		Trace.Fatalln("Failed to open the brand database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	brand, err := getBrand(db, params["brandid"])

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No brand found."))
	}

	fmt.Println(brand)
}

func serverAddBrand(w http.ResponseWriter, r *http.Request) {
	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the brand database."))
		Trace.Fatalln("Failed to open the brand database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	newBrand := Brand{0, "Brand C", "This is brand c"}

	err = addBrand(db, newBrand.Name, newBrand.Description)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error updating brand at database."))
		return
	}

	fmt.Println("Brand successfully added.")
}

func serverEditBrand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the brand database."))
		Trace.Fatalln("Failed to open the brand database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	updatedBrand := Brand{0, "Brand C", "This is brand c"}

	brandID, _ := strconv.Atoi(params["brandid"])

	err = editBrand(db, updatedBrand.Name, updatedBrand.Description, brandID)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error updating brand at database."))
		return
	}
}

func serverDeleteBrand(w http.ResponseWriter, r *http.Request) {
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

	brandID, _ := strconv.Atoi(params["brandid"])

	err = deleteBrand(db, brandID)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error deleting brand from database."))
		return
	}
	fmt.Println("Brand deleted:", brandID)
}

func allCategories(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	categories, err := getCategories(db)
	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error getting categories from database."))
		return
	}
	// json.NewEncoder(w).Encode(productsFromDB)
	fmt.Println(categories)
}

func serverGetCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the categories database."))
		Trace.Fatalln("Failed to open the categories database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	category, err := getCategory(db, params["categoryid"])

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No category found."))
	}

	fmt.Println(category)
}

func serverAddCategory(w http.ResponseWriter, r *http.Request) {
	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the product database."))
		Trace.Fatalln("Failed to open the product database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	newCategory := Category{0, "Category C", "This is category c"}

	err = addCategory(db, newCategory.Name, newCategory.Description)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error adding category at database."))
		return
	}

	fmt.Println("Category successfully added.")
}

func serverEditCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	updatedCategory := Category{0, "Category C", "This is category c"}

	categoryID, _ := strconv.Atoi(params["categoryid"])

	err = editCategory(db, updatedCategory.Name, updatedCategory.Description, categoryID)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error updating category at database."))
		return
	}
}

func serverDeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	categoryID, _ := strconv.Atoi(params["categoryid"])

	err = deleteCategory(db, categoryID)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error deleting category from database."))
		return
	}
	fmt.Println("Category deleted:", categoryID)
}

func main() {

	// main_queries()

	router := mux.NewRouter()
	// router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/products", allproducts)
	router.HandleFunc("/api/v1/products/active", activeProducts)
	router.HandleFunc("/api/v1/products/soldout", soldoutProducts)
	router.HandleFunc("/api/v1/products/unlisted", unlistedProducts)
	router.HandleFunc("/api/v1/product/{productid}", product).Methods("GET")
	router.HandleFunc("/api/v1/product/{productid}", product).Methods("PUT")
	router.HandleFunc("/api/v1/product/{productid}", product).Methods("DELETE")
	router.HandleFunc("/api/v1/product", product).Methods("POST")
	router.HandleFunc("/api/v1/brand", serverAddBrand).Methods("POST")
	router.HandleFunc("/api/v1/brands", allBrands)
	router.HandleFunc("/api/v1/brand/{brandid}", serverGetBrand).Methods("GET")
	router.HandleFunc("/api/v1/brand/{brandid}", serverEditBrand).Methods("PUT")
	router.HandleFunc("/api/v1/brand/{brandid}", serverDeleteBrand).Methods("DELETE")
	router.HandleFunc("/api/v1/category", serverAddCategory).Methods("POST")
	router.HandleFunc("/api/v1/categories", allCategories)
	router.HandleFunc("/api/v1/category/{categoryid}", serverGetCategory).Methods("GET")
	router.HandleFunc("/api/v1/category/{categoryid}", serverEditCategory).Methods("PUT")
	router.HandleFunc("/api/v1/category/{categoryid}", serverDeleteCategory).Methods("DELETE")

	fmt.Println("Listening at port 5000")
	//log.Fatal(http.ListenAndServe(":5000", router))

	err := http.ListenAndServe(":5000", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
