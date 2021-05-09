package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	tpl   *template.Template
	Trace *log.Logger //Prints execution status to stdout, for debugging purposes
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

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

func getActiveProducts(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(products)
	fmt.Println(products)
}

func getSoldoutProducts(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(products)
	fmt.Println(products)
}

func getUnlistedProducts(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(products)
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
	json.NewEncoder(w).Encode(products)
	// fmt.Println(products)
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
			json.NewEncoder(w).Encode(product)
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

		// } else if r.Header.Get("Content-type") == "application/json" {

		//POST is for creating new product
	} else if r.Method == "POST" {
		//read the string sent to the service
		var newProduct Product
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply product information in JSON format"))
			return
		} else {
			//convert JSON to object
			json.Unmarshal(reqBody, &newProduct)

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
		var updatedProduct Product
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply product information in JSON format"))
			return
		} else {
			//convert JSON to object
			json.Unmarshal(reqBody, &updatedProduct)
		}

		//product already exists, update product
		productID, _ := strconv.Atoi(params["productid"])
		err = editProducts(db, updatedProduct.Name, updatedProduct.Image, updatedProduct.DescShort, updatedProduct.DescLong, updatedProduct.Price, updatedProduct.Quantity, updatedProduct.Condition, updatedProduct.CategoryID, updatedProduct.BrandID, updatedProduct.Status, productID)

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
	json.NewEncoder(w).Encode(brands)
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

	json.NewEncoder(w).Encode(brand)

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

	var newBrand Brand
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply product information in JSON format"))
		return
	} else {
		//convert JSON to object
		json.Unmarshal(reqBody, &newBrand)

		//check if product exists; add only if product does not exist
		err = addBrand(db, newBrand.Name, newBrand.Description)

		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error updating product at database."))
			return
		}

		fmt.Println("Brand successfully added.")
	}
}

func serverEditBrand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	brandID, _ := strconv.Atoi(params["brandid"])

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the brand database."))
		Trace.Fatalln("Failed to open the brand database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	var updatedBrand Brand
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply product information in JSON format"))
		return
	} else {
		//convert JSON to object
		json.Unmarshal(reqBody, &updatedBrand)

		err = editBrand(db, updatedBrand.Name, updatedBrand.Description, brandID)

		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error updating brand at database."))
			return
		}

		fmt.Println("Brand successfully added.")
	}
}

func serverDeleteBrand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	brandID, _ := strconv.Atoi(params["brandid"])

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the brand database."))
		Trace.Fatalln("Failed to open the brand database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

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
	json.NewEncoder(w).Encode(categories)
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

	json.NewEncoder(w).Encode(category)

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

	var newCategory Category
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error adding category at database."))
		return
	} else {
		json.Unmarshal(reqBody, &newCategory)

		err := addCategory(db, newCategory.Name, newCategory.Description)

		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error updating product at database."))
			return
		}
	}

	fmt.Println("Category successfully added.")
}

func serverEditCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	catID, _ := strconv.Atoi(params["categoryid"])

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	var updatedCategory Category
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply product information in JSON format"))
		return
	} else {
		//convert JSON to object
		json.Unmarshal(reqBody, &updatedCategory)

		err = editCategory(db, updatedCategory.Name, updatedCategory.Description, catID)

		if err != nil {
			Trace.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Error updating category at database."))
			return
		}
	}
}

func serverDeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	categoryID, _ := strconv.Atoi(params["categoryid"])

	//open the database
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	err = deleteCategory(db, categoryID)

	if err != nil {
		Trace.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Error deleting category from database."))
		return
	}
	fmt.Println("Category deleted:", categoryID)
}

func updateProdQty(w http.ResponseWriter, r *http.Request) {

	// type CartItem struct {
	// 	ID            int
	// 	Name          string
	// 	Price         float64
	// 	QuantityToBuy int
	// }

	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Error opening the category database."))
		Trace.Fatalln("Failed to open the category database.")
	}
	defer db.Close()
	fmt.Println("The database is opened:", db)

	var cartItems []CartItem
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply product information in JSON format"))
		return
	} else {
		json.Unmarshal(reqBody, &cartItems)

		//range through the list of products to be updated
		for _, v := range cartItems {

			productID := strconv.Itoa(v.ID)

			product, err := getProduct(db, productID)

			if err != nil {
				log.Fatalln(err)
			}

			updatedQuantity := product.Quantity - v.QuantityToBuy
			updatedQuantitySold := product.QuantitySold + v.QuantityToBuy

			fmt.Println(updatedQuantity)
			fmt.Println(updatedQuantitySold)

			err = editProductQuantity(db, updatedQuantity, updatedQuantitySold, v.ID)

			if err != nil {
				Trace.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 -Error updating product at database."))
				return
			}
		}
	}

}
