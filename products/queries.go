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

func getBrands(db *sql.DB) ([]Brand, error) {
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

	//fmt.Println(brands)
	return brands, nil
}

func getBrand(db *sql.DB, brandID string) (Brand, error) {
	query := fmt.Sprintf("SELECT * FROM Brands WHERE ID = %s", brandID)

	results, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	var brand Brand

	for results.Next() {
		err = results.Scan(&brand.ID, &brand.Name, &brand.Description, &brand.NumberOfProducts)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}
	}

	//fmt.Println(brands)
	return brand, nil
}

func addBrand(db *sql.DB, Name string, Description string, NumOfProducts int) error {
	query := fmt.Sprintf("INSERT INTO Brands (Name, Description, Number_Of_Products) VALUES ('%s', '%s', %d)", Name, Description, NumOfProducts)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func editBrand(db *sql.DB, Name string, Description string, NumOfProducts int, ID int) error {
	query := fmt.Sprintf("UPDATE Brands SET Name='%s', Description='%s', Number_Of_Products=%d WHERE ID=%d", Name, Description, NumOfProducts, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func deleteBrand(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Brands WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func getCategories(db *sql.DB) ([]Category, error) {
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

	//fmt.Println(categories)
	return categories, nil
}

func getCategory(db *sql.DB, categoryID string) (Category, error) {
	query := fmt.Sprintf("SELECT * FROM Categories WHERE ID = %s", categoryID)

	results, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	var category Category

	for results.Next() {
		err = results.Scan(&category.ID, &category.Name, &category.Description, &category.NumberOfProducts)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}
	}

	//fmt.Println(brands)
	return category, nil
}

func addCategory(db *sql.DB, Name string, Description string, NumOfProducts int) error {
	query := fmt.Sprintf("INSERT INTO Categories (Name, Description, Number_Of_Products) VALUES ('%s', '%s', %d)", Name, Description, NumOfProducts)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func editCategory(db *sql.DB, Name string, Description string, NumOfProducts int, ID int) error {
	query := fmt.Sprintf("UPDATE Categories SET Name='%s', Description='%s', Number_Of_Products=%d WHERE ID=%d", Name, Description, NumOfProducts, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func deleteCategory(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Categories WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func getProducts(db *sql.DB) ([]Product, error) {
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

	//fmt.Println(products)
	return products, nil
}

func getProduct(db *sql.DB, productID string) ([]Product, error) {
	query := fmt.Sprintf("SELECT * FROM Products WHERE ID='%s'", productID)
	results, err := db.Query(query)

	if err != nil {
		fmt.Println("queries: line 193: ", err)
		log.Fatalln()
	}

	var products []Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.Details, &product.DateAdded, &product.Price, &product.Quantity, &product.CategoryID, &product.BrandID)

		if err != nil {
			fmt.Println("queries: line 205: ", err)
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	//fmt.Println(products)
	return products, nil
}

func addProducts(db *sql.DB, Name string, Image string, Details string, DateAdded string, Price float64, Quantity int, CategoryID int, BrandID int) error {
	query := fmt.Sprintf("INSERT INTO Products (Name, Image, Details, Date_Added, Price, Quantity, Category_ID, Brand_ID) VALUES ('%s', '%s', '%s', '%s', %v, %d, %d, %d)", Name, Image, Details, DateAdded, Price, Quantity, CategoryID, BrandID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func editProducts(db *sql.DB, Name string, Image string, Details string, DateAdded string, Price float64, Quantity int, CategoryID int, BrandID int, ID int) error {
	query := fmt.Sprintf("UPDATE Products SET Name='%s', Image='%s', Details='%s', Date_Added='%s', Price=%.2f, Quantity=%d, Category_ID=%d, Brand_ID=%d WHERE ID=%d", Name, Image, Details, DateAdded, Price, Quantity, CategoryID, BrandID, ID)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func deleteProducts(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Products WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func main_queries() {
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

	addBrand(db, "Brand A2", "This is Brand A2", 130)
	addCategory(db, "Category A0002", "This is category A0001", 30)
	addProducts(db, "Bag", "nil", "This is a bag", "2021-04-27", 20.50, 5, 2, 2)

	// editProducts(db, "Bag", "nil", "This is a bag", "2021-04-27", 20.50, 5, 1, 1, 1)
	// deleteProducts(db, 1)
	// getProducts(db)

	// editBrand(db, "Brand A1", "This is Brand A1", 10, 1)
	// deleteBrand(db, 1)
	// getBrands(db)

	// editCategory(db, "Category A0001", "This is category A0001", 100, 1)
	// deleteCategory(db, 1)
	// getCategories(db)

}
