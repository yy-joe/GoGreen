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
	ID           int
	Name         string
	Image        string
	DescShort    string
	DescLong     string
	DateCreated  string
	DateModified string
	Price        float64
	Quantity     int
	Condition    string
	CategoryID   int
	BrandID      int
	Status       string
}

type Brand struct {
	ID          int
	Name        string
	Description string
}

type Category struct {
	ID          int
	Name        string
	Description string
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
		err = results.Scan(&brand.ID, &brand.Name, &brand.Description)

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
		err = results.Scan(&brand.ID, &brand.Name, &brand.Description)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}
	}

	//fmt.Println(brands)
	return brand, nil
}

func addBrand(db *sql.DB, Name string, Description string) error {
	query := fmt.Sprintf("INSERT INTO Brands (Name, Description) VALUES ('%s', '%s')", Name, Description)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func editBrand(db *sql.DB, Name string, Description string, ID int) error {
	query := fmt.Sprintf("UPDATE Brands SET Name='%s', Description='%s' WHERE ID=%d", Name, Description, ID)

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

		err = results.Scan(&category.ID, &category.Name, &category.Description)

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
		err = results.Scan(&category.ID, &category.Name, &category.Description)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}
	}

	//fmt.Println(brands)
	return category, nil
}

func addCategory(db *sql.DB, Name string, Description string) error {
	query := fmt.Sprintf("INSERT INTO Categories (Name, Description) VALUES ('%s', '%s')", Name, Description)

	_, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	return nil
}

func editCategory(db *sql.DB, Name string, Description string, ID int) error {
	query := fmt.Sprintf("UPDATE Categories SET Name='%s', Description='%s' WHERE ID=%d", Name, Description, ID)

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

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

		if err != nil {
			fmt.Println(err)
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	//fmt.Println(products)
	return products, nil
}

func getProductsByStatus(db *sql.DB, status string) ([]Product, error) {
	query := fmt.Sprintf("SELECT * FROM GoGreen.Products WHERE Status='%s'", status)
	results, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		log.Fatalln()
	}

	var products []Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

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

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

		if err != nil {
			fmt.Println("queries: line 205: ", err)
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	//fmt.Println(products)
	return products, nil
}

func addProducts(db *sql.DB, Name string, Image string, DescShort string, DescLong string, Price float64, Quantity int, Condition string, CategoryID int, BrandID int, Status string) error {
	query := fmt.Sprintf("INSERT INTO Products (Name, Image, Desc_Short, Desc_Long, Date_Created, Date_Modified, Price, Quantity, `Condition`, Category_ID, Brand_ID, Status) VALUES ('%s', '%s', '%s', '%s', curdate(), curdate(), %v, %d, '%s', %d, %d, '%s')", Name, Image, DescShort, DescLong, Price, Quantity, Condition, CategoryID, BrandID, Status)

	_, err := db.Query(query)

	if err != nil {
		// fmt.Println(err)
		log.Fatalln(err)
	}

	return nil
}

func editProducts(db *sql.DB, Name string, Image string, DescShort string, DescLong string, Price float64, Quantity int, Condition string, CategoryID int, BrandID int, Status string, ID int) error {
	query := fmt.Sprintf("UPDATE Products SET Name='%s', Image='%s', Desc_Short='%s', Desc_Long='%s', Date_Modified=curdate(), Price=%.2f, Quantity=%d, `Condition`='%s', Category_ID=%d, Brand_ID=%d, Status='%s' WHERE ID=%d", Name, Image, DescShort, DescLong, Price, Quantity, Condition, CategoryID, BrandID, Status, ID)

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

	// addBrand(db, "Brand A", "This is Brand A", 0)
	// addCategory(db, "Category A", "This is category A", 0)
	addProducts(db, "Bag", "nil", "This is a bag", "This is a very very very big bag.", 20.50, 5, "New", 1, 1, "Live")

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