package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// getProducts, editProduct, addProduct, deleteProduct
// getBrand, editBrand, addBrand, deleteBrand
// getCategory, editCategory, addCategory, deleteCategory

type Product struct {
	ID         int
	Name       string
	Image      string
	Details    string
	DateAdded  string
	Price      float64
	Quantity   int
	CategoryID int
	BrandID    int
}

type Brand struct {
	ID               int
	Name             string
	Description      string
	NumberOfProducts int
}

type Category struct {
	ID               int
	Name             string
	Description      string
	NumberOfProducts int
}

func getBrands(db *sql.DB) {
	results, err := db.Query("SELECT * FROM GoGreen.Brands")

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	var brands []Brand

	for results.Next() {
		var brand Brand
		err = results.Scan(&brand.ID, &brand.Name, &brand.Description, &brand.NumberOfProducts)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}

		brands = append(brands, brand)
	}

	fmt.Println(brands)
}

func addBrand(db *sql.DB, Name string, Description string, NumOfProducts int) {
	query := fmt.Sprintf("INSERT INTO Brands (Name, Description, Number_Of_Products) VALUES ('%s', '%s', %d)", Name, Description, NumOfProducts)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func getCategories(db *sql.DB) {
	results, err := db.Query("SELECT * FROM GoGreen.Categories")

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	var categories []Category

	for results.Next() {
		var category Category

		err = results.Scan(&category.ID, &category.Name, &category.Description, &category.NumberOfProducts)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}

		categories = append(categories, category)
	}

	fmt.Println(categories)
}

func addCategory(db *sql.DB, Name string, Description string, NumOfProducts int) {
	query := fmt.Sprintf("INSERT INTO Categories (Name, Description, Number_Of_Products) VALUES ('%s', '%s', %d)", Name, Description, NumOfProducts)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func getProducts(db *sql.DB) {
	results, err := db.Query("SELECT * FROM GoGreen.Products")

	if err != nil {
		fmt.Println(err)
		log.Fatalln()
	}

	var products []Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.Details, &product.DateAdded, &product.Price, &product.Quantity, &product.CategoryID, &product.BrandID)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	fmt.Println(products)
}

func addProducts(db *sql.DB, ID int, Name string, Image string, Details string, DateAdded string, Price float64, Quantity int, CategoryID int, BrandID int) {
	query := fmt.Sprintf("INSERT INTO Products VALUES (%d, '%s', '%s', '%s', '%s', %v, %d, %d, %d)", ID, Name, Image, Details, DateAdded, Price, Quantity, CategoryID, BrandID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func main() {
	//Use mysql as driverName and a valid DSN as data source name
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/GoGreen")

	//handle error
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	//defer the close till after the main function has finished executing
	defer db.Close()

	fmt.Println("Database opened")

	// addBrand(db, "Brand A", "This is Brand A", 0)
	getBrands(db)

	// addCategory(db, "Category A", "This is category A", 0)
	getCategories(db)

	// addProducts(db, 1, "Bag", "nil", "This is a bag", "2021-04-27", 20.50, 5, 1, 1)
	getProducts(db)
}
