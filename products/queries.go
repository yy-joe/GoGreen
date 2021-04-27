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

/*
func EditRecord(db *sql.DB, ID int, FN string, LN string, Age int) {
	query := fmt.Sprintf("UPDATE Persons SET FirstName='%s', LastName='%s', Age=%d WHERE ID=%d", FN, LN, Age, ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
*/

func editBrand(db *sql.DB, Name string, Description string, NumOfProducts int, ID int) {
	query := fmt.Sprintf("UPDATE Brands SET Name='%s', Description='%s', Number_Of_Products=%d WHERE ID=%d", Name, Description, NumOfProducts, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func deleteBrand(db *sql.DB, ID int) {
	query := fmt.Sprintf("DELETE FROM Brands WHERE ID='%d'", ID)
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

func editCategory(db *sql.DB, Name string, Description string, NumOfProducts int, ID int) {
	query := fmt.Sprintf("UPDATE Categories SET Name='%s', Description='%s', Number_Of_Products=%d WHERE ID=%d", Name, Description, NumOfProducts, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func deleteCategory(db *sql.DB, ID int) {
	query := fmt.Sprintf("DELETE FROM Categories WHERE ID='%d'", ID)
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

func editProducts(db *sql.DB, Name string, Image string, Details string, DateAdded string, Price float64, Quantity int, CategoryID int, BrandID int, ID int) {
	query := fmt.Sprintf("UPDATE Products SET Name='%s', Image='%s', Details='%s', Date_Added='%s', Price=%.2f, Quantity=%d, Category_ID=%d, Brand_ID=%d WHERE ID=%d", Name, Image, Details, DateAdded, Price, Quantity, CategoryID, BrandID, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func deleteProducts(db *sql.DB, ID int) {
	query := fmt.Sprintf("DELETE FROM Products WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func main() {
	//Use mysql as driverName and a valid DSN as data source name
	db, err := sql.Open("mysql", "root:QQ2kepiting@tcp(127.0.0.1:3306)/GoGreen")

	//handle error
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	//defer the close till after the main function has finished executing
	defer db.Close()

	fmt.Println("Database opened")

	// // editProducts(db, "Bag", "nil", "This is a bag", "2021-04-27", 20.50, 5, 1, 1, 1)
	// deleteProducts(db, 1)
	// getProducts(db)

	// // editBrand(db, "Brand A1", "This is Brand A1", 10, 1)
	// deleteBrand(db, 1)
	// getBrands(db)

	// editCategory(db, "Category A0001", "This is category A0001", 100, 1)
	deleteCategory(db, 1)
	getCategories(db)

}
