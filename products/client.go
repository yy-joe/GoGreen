package products

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// transport layer security Configuration
var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

var mutex sync.Mutex

const baseURL = "https://localhost:3000/api/v1/admin/"

func ProdMain(w http.ResponseWriter, r *http.Request) {
	// if len(storedProducts) == 0 {
	// 	err := initGlobalVars()
	// 	if err != nil {
	// 		errData := ErrorTplData{
	// 			ErrorMessage: "Error getting data from the server.",
	// 			RedirectLink: "/products/all",
	// 			ButtonValue:  "Try again",
	// 		}
	// 		tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
	// 		return
	// 	}
	// }
	//check if global vars are initialized.

	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "The product list cannot be loaded at the moment.",
				RedirectLink: "/products/all",
				ButtonValue:  "Try again",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	url := baseURL + "products"
	fmt.Println(url)
	// res, err := http.Get(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		panic(err)
	}

	var products []Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		panic(err)
	}

	switch sortKey {
	case "":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", formatProdTplData(products))
	case "Price":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", formatProdTplData(productsByPrice))

	case "Quantity":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", formatProdTplData(productsByQuantity))

	default:
		sortedProducts := sortProducts(products, sortKey)
		fmt.Println(sortedProducts)

		tpl.ExecuteTemplate(w, "prodMain.gohtml", formatProdTplData(sortedProducts))
	}
}

type prodTpl struct {
	ID           int
	Name         string
	Image        string
	DescShort    string
	DescLong     string
	DateCreated  string
	DateModified string
	Price        float64
	Quantity     int
	QuantitySold int
	Condition    string
	CategoryName string
	BrandName    string
	Status       string
}

func formatProdTplData(products []Product) (productsTpl []prodTpl) {
	storedCatMap := make(map[int]string)
	storedBrandMap := make(map[int]string)

	//populate the maps
	for _, v := range storedCategories {
		storedCatMap[v.ID] = v.Name
	}
	for _, v := range storedBrands {
		storedBrandMap[v.ID] = v.Name
	}
	fmt.Println("------storedCatMap = ", storedCatMap)
	fmt.Println("------storedBrandMap = ", storedBrandMap)

	for _, v := range products {
		var prod prodTpl
		prod.ID = v.ID
		prod.Name = v.Name
		prod.Image = v.Image
		prod.DescShort = v.DescShort
		prod.DescLong = v.DescLong
		prod.DateCreated = v.DateCreated
		prod.DateModified = v.DateModified
		prod.Price = v.Price
		prod.Quantity = v.Quantity
		prod.QuantitySold = v.QuantitySold
		prod.Condition = v.Condition
		prod.CategoryName = storedCatMap[v.CategoryID]
		prod.BrandName = storedBrandMap[v.BrandID]
		prod.Status = v.Status
		productsTpl = append(productsTpl, prod)
	}
	fmt.Println("------productsTpl = ", productsTpl)
	return
}

func sortProducts(products []Product, sortKey string) []Product {
	return mergeSort(products, sortKey)
	//return selectionSort(products, sortKey)
}

func ProdByStatus(w http.ResponseWriter, r *http.Request) {
	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	params := mux.Vars(r)
	byStatus := params["byStatus"]

	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "The product list cannot be loaded at the moment.",
				RedirectLink: "/products/" + byStatus,
				ButtonValue:  "Try again",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	url := baseURL + "products/" + byStatus
	fmt.Println(url)
	// res, err := http.Get(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		panic(err)
	}

	var products []Product
	err = json.Unmarshal(data, &products)

	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		panic(err)
	}

	// productsByPrice = make([]Product, 0, len(products))
	// productsByQuantity = make([]Product, 0, len(products))
	// productsByPrice = mergeSort(products, "Price")
	// productsByQuantity = mergeSort(products, "Quantity")

	Status := strings.Title(byStatus)
	type templateData struct {
		Status   string
		ByStatus string
		Products []Product
	}

	// switch sortKey {
	// case "":
	// 	tpl.ExecuteTemplate(w, "prodByStatus.gohtml", templateData{Status, byStatus, products})
	// case "Price":
	// 	tpl.ExecuteTemplate(w, "prodByStatus.gohtml", templateData{Status, byStatus, productsByPrice})

	// case "Quantity":
	// 	tpl.ExecuteTemplate(w, "prodByStatus.gohtml", templateData{Status, byStatus, productsByQuantity})

	// default:
	if sortKey == "" {
		sortKey = "Name"
	}
	sortedProducts := sortProducts(products, sortKey)
	fmt.Println(sortedProducts)
	fmt.Println("Status=", Status)
	fmt.Println("byStatus=", byStatus)

	tpl.ExecuteTemplate(w, "prodByStatus.gohtml", templateData{Status, byStatus, sortedProducts})
	// }
}

func ProdAdd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "Failed to add new product.",
				RedirectLink: "/products/all",
				ButtonValue:  "Back to main product page",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	if r.Method == http.MethodPost {
		fmt.Println("prodAdd: processing submitted form...")
		//add the new product to the database
		price, _ := strconv.ParseFloat(r.FormValue("Price"), 64)
		qty, _ := strconv.Atoi(r.FormValue("Quantity"))
		catid, _ := strconv.Atoi(r.FormValue("CategoryID"))
		brandid, _ := strconv.Atoi(r.FormValue("BrandID"))
		newProduct := Product{
			ID:           0,
			Name:         r.FormValue("Name"),
			Image:        "nil",
			DescShort:    r.FormValue("DescShort"),
			DescLong:     r.FormValue("DescLong"),
			DateCreated:  "",
			DateModified: "",
			Price:        price,
			Quantity:     qty,
			Condition:    r.FormValue("Condition"),
			CategoryID:   catid,
			BrandID:      brandid,
			Status:       r.FormValue("Status"),
		}

		// json.NewEncoder(w).Encode(newProduct)
		jsonValue, err := json.Marshal(newProduct)
		if err != nil {
			fmt.Printf("Error marshaling json data: %s\n", err)
			panic(err)
		}

		url := baseURL + "product"
		fmt.Println(url)
		// _, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		//update global vars
		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}
		json.Unmarshal(reqBody, &newProduct)
		mutex.Lock()
		{
			storedProducts = append(storedProducts, newProduct)
		}
		mutex.Unlock()

		//update sorted product list
		mutex.Lock()
		{
			productsByPrice = mergeSort(storedProducts, "Price")
			productsByQuantity = mergeSort(storedProducts, "Quantity")
		}
		mutex.Unlock()

		// fmt.Println("Updated storedProducts :", storedProducts)
		//direct user back to the main products page
		http.Redirect(w, r, "/products/all", http.StatusSeeOther)
		return
	}

	//get the categories & brands
	catsAndBrands := struct {
		Categories []Category
		Brands     []Brand
	}{
		Categories: storedCategories,
		Brands:     storedBrands,
	}
	// fmt.Println(catsAndBrands)
	tpl.ExecuteTemplate(w, "prodAdd.gohtml", catsAndBrands)
}

func ProdUpdate(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "Failed to update product.",
				RedirectLink: "/products/all",
				ButtonValue:  "Back to main product page",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	params := mux.Vars(r)
	productID := params["productid"]

	productIntID, _ := strconv.Atoi(productID)

	if r.Method == http.MethodPost {
		//update the product at database
		price, _ := strconv.ParseFloat(r.FormValue("Price"), 64)
		qty, _ := strconv.Atoi(r.FormValue("Quantity"))
		catid, _ := strconv.Atoi(r.FormValue("CategoryID"))
		brandid, _ := strconv.Atoi(r.FormValue("BrandID"))
		updatedProduct := Product{
			ID:           0,
			Name:         r.FormValue("Name"),
			Image:        "nil",
			DescShort:    r.FormValue("DescShort"),
			DescLong:     r.FormValue("DescLong"),
			DateCreated:  "",
			DateModified: "",
			Price:        price,
			Quantity:     qty,
			QuantitySold: 0,
			Condition:    r.FormValue("Condition"),
			CategoryID:   catid,
			BrandID:      brandid,
			Status:       r.FormValue("Status"),
		}

		// json.NewEncoder(w).Encode(newProduct)
		jsonValue, err := json.Marshal(updatedProduct)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %s\n", err)
			panic(err)
		}

		url := baseURL + "product/" + productID

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		if res.StatusCode != 200 {
			return
		}

		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}
		json.Unmarshal(reqBody, &updatedProduct)

		for i, v := range storedProducts {
			if v.ID == productIntID {
				mutex.Lock()
				{
					storedProducts[i].Name = updatedProduct.Name
					storedProducts[i].Image = updatedProduct.Image
					storedProducts[i].DescShort = updatedProduct.DescShort
					storedProducts[i].DescLong = updatedProduct.DescLong
					storedProducts[i].DateModified = updatedProduct.DateModified
					storedProducts[i].Price = updatedProduct.Price
					storedProducts[i].Quantity = updatedProduct.Quantity
					storedProducts[i].Condition = updatedProduct.Condition
					storedProducts[i].CategoryID = updatedProduct.CategoryID
					storedProducts[i].BrandID = updatedProduct.BrandID
					storedProducts[i].Status = updatedProduct.Status
				}
				mutex.Unlock()
				break
			}
		}

		// fmt.Println("!!!! Updated storedProducts :", storedProducts)

		//direct user back to the main products page
		byStatus := r.FormValue("byStatus")
		if byStatus == "" {
			http.Redirect(w, r, "/products/all", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/products/"+byStatus, http.StatusSeeOther)
		}

	}
}

func ProdDetail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "Failed to get product detail.",
				RedirectLink: "/products/all",
				ButtonValue:  "Back to main product page",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	id := params["productid"]

	url := baseURL + "product/" + id
	fmt.Println(url)
	// res, err := http.Get(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body %s\n", err)
		panic(err)
	}

	var product Product
	err = json.Unmarshal(data, &product)

	if err != nil {
		fmt.Printf("Error reading response body %s\n", err)
		panic(err)
	}

	byStatus := r.FormValue("byStatus")
	var showByStatus bool
	if byStatus != "" {
		showByStatus = true
	}

	templateData := struct {
		Categories   []Category
		Brands       []Brand
		Product      Product
		ShowByStatus bool
		ByStatus     string
	}{
		Categories:   storedCategories,
		Brands:       storedBrands,
		Product:      product,
		ShowByStatus: showByStatus,
		ByStatus:     byStatus,
	}

	tpl.ExecuteTemplate(w, "prodDetail.gohtml", templateData)
}

func ProdDelete(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//reload the shopping cart template
			errData := ErrorTplData{
				ErrorMessage: "Failed deleting product.",
				RedirectLink: "/products/all",
				ButtonValue:  "Back to main product page",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	productID := params["productid"]
	prodIntID, _ := strconv.Atoi(productID)

	url := baseURL + "product/" + productID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	if res.StatusCode != 200 {
		return
	}

	//update global vars
	for i, v := range storedProducts {
		if v.ID == prodIntID {
			mutex.Lock()
			{
				storedProducts = append(storedProducts[:i], storedProducts[i+1:]...)
			}
			mutex.Unlock()
			break
		}
	}
	fmt.Println("Deleted a product : ", storedProducts)

	//direct user back to the main products page
	http.Redirect(w, r, "/products/all", http.StatusSeeOther)
}

func clientGetCategories() (categories []Category, err error) {
	url := baseURL + "categories"
	fmt.Println(url)
	// res, err := http.Get(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body %s\n", err)
		return
	}

	err = json.Unmarshal(data, &categories)
	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		return
	}
	return
}

func CatMain(w http.ResponseWriter, r *http.Request) {
	if len(storedProducts) == 0 {
		initGlobalVars()
	}
	categories, err := clientGetCategories()
	if err != nil {
		errData := ErrorTplData{
			ErrorMessage: "Error getting category data from the server.",
			RedirectLink: "/categories/all",
			ButtonValue:  "Try again",
		}
		tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		return
	}
	tpl.ExecuteTemplate(w, "catMain.gohtml", categories)
}

func CatDetail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed getting category detail.",
				RedirectLink: "/categories/all",
				ButtonValue:  "Back to categories",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()
	params := mux.Vars(r)
	catID := params["categoryid"]

	url := baseURL + "category/" + catID
	fmt.Println(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		panic(err)
	}

	var category Category
	err = json.Unmarshal(data, &category)

	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		panic(err)
	}

	tpl.ExecuteTemplate(w, "catDetail.gohtml", category)
}

func CatAdd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed adding new category.",
				RedirectLink: "/categories/all",
				ButtonValue:  "Back to categories",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	if r.Method == http.MethodPost {

		newCategory := Category{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(newCategory)
		if err != nil {
			fmt.Printf("Error marshaling JSON data: %s\n", err)
			panic(err)
		}

		url := baseURL + "category"

		res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		//update the global variable storedCategories
		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}
		json.Unmarshal(reqBody, &newCategory)
		mutex.Lock()
		{
			storedCategories = append(storedCategories, newCategory)
		}
		mutex.Unlock()
		// fmt.Println("Updated storedCategories :", storedCategories)

		http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "catAdd.gohtml", nil)
}

func CatUpdate(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed updating category.",
				RedirectLink: "/categories/all",
				ButtonValue:  "Back to categories",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	catID := params["categoryid"]

	catIntID, _ := strconv.Atoi(catID)

	if r.Method == http.MethodPost {
		updatedCategory := Category{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(updatedCategory)
		if err != nil {
			fmt.Printf("Error marshaling JSON data: %s\n", err)
			panic(err)
		}

		url := baseURL + "category/" + catID
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		if res.StatusCode != 200 {
			return
		}

		reqBody, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}
		json.Unmarshal(reqBody, &updatedCategory)

		for i, v := range storedCategories {
			if v.ID == catIntID {
				mutex.Lock()
				{
					storedCategories[i].Name = updatedCategory.Name
					storedCategories[i].Description = updatedCategory.Description
				}
				mutex.Unlock()
				break
			}
		}

		// fmt.Println("!!!! Updated storedCategories :", storedCategories)

		//direct user back to the main products page
		http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
	}
}

func CatDelete(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed deleting category.",
				RedirectLink: "/categories/all",
				ButtonValue:  "Back to categories",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	catID := params["categoryid"]
	catIntID, _ := strconv.Atoi(catID)

	url := baseURL + "category/" + catID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	if res.StatusCode != 200 {
		return
	}

	//update global vars
	for i, v := range storedCategories {
		if v.ID == catIntID {
			mutex.Lock()
			{
				storedCategories = append(storedCategories[:i], storedCategories[i+1:]...)
			}
			mutex.Unlock()
			break
		}
	}
	// fmt.Println("Deleted a category : ", storedCategories)

	//direct user back to the main products page
	http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
}

func clientGetProducts() (products []Product, err error) {

	url := baseURL + "products"
	fmt.Println(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	err = json.Unmarshal(data, &products)
	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		return
	}
	return
}

func clientGetBrands() (brands []Brand, err error) {

	url := baseURL + "brands"
	fmt.Println(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	err = json.Unmarshal(data, &brands)
	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		return
	}
	return
}

func BrandMain(w http.ResponseWriter, r *http.Request) {
	if len(storedProducts) == 0 {
		initGlobalVars()
	}
	tpl.ExecuteTemplate(w, "brandMain.gohtml", storedBrands)
}

func BrandDetail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed getting brand details.",
				RedirectLink: "/brands/all",
				ButtonValue:  "Back to brands",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	brandID := params["brandid"]

	url := baseURL + "brand/" + brandID
	fmt.Println(url)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		panic(err)
	}

	var brand Brand
	err = json.Unmarshal(data, &brand)
	if err != nil {
		fmt.Printf("Error reading JSON data: %s\n", err)
		panic(err)
	}

	tpl.ExecuteTemplate(w, "brandDetail.gohtml", brand)
}

func BrandAdd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed adding brand.",
				RedirectLink: "/brands/all",
				ButtonValue:  "Back to brands",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	if r.Method == http.MethodPost {

		newBrand := Brand{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(newBrand)
		if err != nil {
			fmt.Printf("Error marshaling JSON data: %s\n", err)
			panic(err)
		}

		url := baseURL + "brand"

		res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		//update the global variable storedBrands
		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}
		json.Unmarshal(reqBody, &newBrand)
		mutex.Lock()
		{
			storedBrands = append(storedBrands, newBrand)
		}
		mutex.Unlock()
		// fmt.Println("Updated storedBrands :", storedBrands)

		http.Redirect(w, r, "/brands/all", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "brandAdd.gohtml", nil)
}

func BrandUpdate(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed updating brand.",
				RedirectLink: "/brands/all",
				ButtonValue:  "Back to brands",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	brandID := params["brandid"]

	brandIntID, _ := strconv.Atoi(brandID)

	if r.Method == http.MethodPost {
		updatedBrand := Brand{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(updatedBrand)
		if err != nil {
			fmt.Printf("Error marshaling JSON data: %s\n", err)
			panic(err)
		}

		url := baseURL + "brand/" + brandID

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			panic(err)
		}

		if res.StatusCode != 200 {
			return
		}

		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			panic(err)
		}

		json.Unmarshal(reqBody, &updatedBrand)

		for i, v := range storedBrands {
			if v.ID == brandIntID {
				mutex.Lock()
				{
					storedBrands[i].Name = updatedBrand.Name
					storedBrands[i].Description = updatedBrand.Description
				}
				mutex.Unlock()
				break
			}
		}

		// fmt.Println("!!!! Updated storedBrands :", storedBrands)

		//direct user back to the main products page
		http.Redirect(w, r, "/brands/all", http.StatusSeeOther)
	}
}

func BrandDelete(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errData := ErrorTplData{
				ErrorMessage: "Failed deleting brand.",
				RedirectLink: "/brands/all",
				ButtonValue:  "Back to brands",
			}
			tpl.ExecuteTemplate(w, "errorPage.gohtml", errData)
		}
	}()

	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	params := mux.Vars(r)
	brandID := params["brandid"]
	brandIntID, _ := strconv.Atoi(brandID)

	url := baseURL + "brand/" + brandID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	if res.StatusCode != 200 {
		return
	}

	//update global vars
	for i, v := range storedBrands {
		if v.ID == brandIntID {
			mutex.Lock()
			{
				storedBrands = append(storedBrands[:i], storedBrands[i+1:]...)
			}
			mutex.Unlock()
			break
		}
	}
	// fmt.Println("Deleted a brand : ", storedBrands)

	//direct user back to the main products page
	http.Redirect(w, r, "/brands/all", http.StatusSeeOther)
}
