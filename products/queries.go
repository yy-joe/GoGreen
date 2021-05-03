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
		log.Fatalln(err)
	}

	var brands []Brand

	for results.Next() {
		var brand Brand
		err = results.Scan(&brand.ID, &brand.Name, &brand.Description)

		if err != nil {
			log.Fatalln(err)
		}

		brands = append(brands, brand)
	}
	return brands, err
}

func getBrand(db *sql.DB, brandID string) (Brand, error) {
	var brand Brand

	err := db.QueryRow("SELECT * FROM Brands WHERE ID=?", brandID).Scan(&brand.ID, &brand.Name, &brand.Description)

	if err != nil {
		log.Fatalln(err)
	}
	return brand, err
}

func addBrand(db *sql.DB, Name string, Description string) error {
	query := fmt.Sprintf("INSERT INTO Brands (Name, Description) VALUES ('%s', '%s')", Name, Description)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func editBrand(db *sql.DB, Name string, Description string, ID int) error {
	query := fmt.Sprintf("UPDATE Brands SET Name='%s', Description='%s' WHERE ID=%d", Name, Description, ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func deleteBrand(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Brands WHERE ID='%d'", ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func getCategories(db *sql.DB) ([]Category, error) {
	results, err := db.Query("SELECT * FROM GoGreen.Categories")

	if err != nil {
		log.Fatalln(err)
	}

	var categories []Category

	for results.Next() {
		var category Category

		err = results.Scan(&category.ID, &category.Name, &category.Description)

		if err != nil {
			log.Fatalln(err)
		}

		categories = append(categories, category)
	}
	return categories, err
}

func getCategory(db *sql.DB, categoryID string) (Category, error) {
	var category Category

	err := db.QueryRow("SELECT * FROM Categories WHERE ID=?", categoryID).Scan(&category.ID, &category.Name, &category.Description)

	if err != nil {
		log.Fatalln(err)
	}

	return category, err
}

func addCategory(db *sql.DB, Name string, Description string) error {
	query := fmt.Sprintf("INSERT INTO Categories (Name, Description) VALUES ('%s', '%s')", Name, Description)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func editCategory(db *sql.DB, Name string, Description string, ID int) error {
	query := fmt.Sprintf("UPDATE Categories SET Name='%s', Description='%s' WHERE ID=%d", Name, Description, ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func deleteCategory(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Categories WHERE ID='%d'", ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func getProducts(db *sql.DB) ([]Product, error) {
	results, err := db.Query("SELECT * FROM GoGreen.Products")

	if err != nil {
		log.Fatalln(err)
	}

	var products []Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

		if err != nil {
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	return products, err
}

func getProductsByStatus(db *sql.DB, status string) ([]Product, error) {
	query := fmt.Sprintf("SELECT * FROM GoGreen.Products WHERE Status='%s'", status)

	results, err := db.Query(query)

	if err != nil {
		log.Fatalln(err)
	}

	var products []Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

		if err != nil {
			log.Fatalln(err)
		}

		products = append(products, product)
	}

	return products, err
}

func getProduct(db *sql.DB, productID string) (Product, error) {
	var product Product

	err := db.QueryRow("SELECT * FROM Products WHERE ID=?", productID).Scan(&product.ID, &product.Name, &product.Image, &product.DescShort, &product.DescLong, &product.DateCreated, &product.DateModified, &product.Price, &product.Quantity, &product.Condition, &product.CategoryID, &product.BrandID, &product.Status)

	if err != nil {
		log.Fatalln(err)
	}

	return product, err
}

func addProducts(db *sql.DB, Name string, Image string, DescShort string, DescLong string, Price float64, Quantity int, Condition string, CategoryID int, BrandID int, Status string) error {
	query := fmt.Sprintf("INSERT INTO Products (Name, Image, Desc_Short, Desc_Long, Date_Created, Date_Modified, Price, Quantity, `Condition`, Category_ID, Brand_ID, Status) VALUES ('%s', '%s', '%s', '%s', curdate(), curdate(), %v, %d, '%s', %d, %d, '%s')", Name, Image, DescShort, DescLong, Price, Quantity, Condition, CategoryID, BrandID, Status)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func editProducts(db *sql.DB, Name string, Image string, DescShort string, DescLong string, Price float64, Quantity int, Condition string, CategoryID int, BrandID int, Status string, ID int) error {
	query := fmt.Sprintf("UPDATE Products SET Name='%s', Image='%s', Desc_Short='%s', Desc_Long='%s', Date_Modified=curdate(), Price=%.2f, Quantity=%d, `Condition`='%s', Category_ID=%d, Brand_ID=%d, Status='%s' WHERE ID=%d", Name, Image, DescShort, DescLong, Price, Quantity, Condition, CategoryID, BrandID, Status, ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func deleteProducts(db *sql.DB, ID int) error {
	query := fmt.Sprintf("DELETE FROM Products WHERE ID='%d'", ID)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln(err)
	}

	return err
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

	// addBrand(db, "Brand A", "This is Brand A")
	// addCategory(db, "Category A", "This is category A")
	// addProducts(db, "Bag", "nil", "This is a bag", "This is a very very very big bag.", 20.50, 5, "New", 1, 1, "active")

	// //add a few more brands
	// addBrand(db, "Brand B", "This is Brand B")
	// addBrand(db, "Brand C", "This is Brand C")

	// //add a few more categories
	// addCategory(db, "Category B", "This is category B")
	// addCategory(db, "Category C", "This is category C")

	// //add a few more products
	// addProducts(db, "Bag B", "nil", "This is another bag", "This is also a very big bag.", 16, 10, "New", 1, 2, "active")
	// addProducts(db, "Bag C", "nil", "This is one more bag", "This is not a very big bag.", 8.80, 5, "New", 1, 1, "unlisted")
	// addProducts(db, "Lunch bag 1", "nil", "This is an insulated lunch bag", "This insulated lunch bag is ideal to keep your food warm/cool.", 20, 10, "New", 2, 2, "active")
	// addProducts(db, "CutleryXYZ", "nil", "This is a set of reusable cutleries", "The package contains a spoon, a fork, a knife and a pair of chopsticks.", 8, 10, "New", 3, 3, "soldout")
	// addProducts(db, "MyStraw", "nil", "Metal straw", "This is a reusable straw.", 5, 10, "New", 3, 2, "active")
	// addProducts(db, "Foodbox", "nil", "This is a lunch box.", "This medium sized lunch box is big enough to store takeaway food, yet small enough to carry.", 25, 8, "New", 3, 1, "unlisted")
	// addProducts(db, "Baggy", "nil", "Large shopping bag", "This is a huge shopping bag.", 20, 10, "New", 1, 3, "active")

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
