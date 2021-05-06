package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	productsByPrice    []Product
	productsByQuantity []Product
)

const baseURL = "http://localhost:5000/api/v1/admin/"

func prodMain(w http.ResponseWriter, r *http.Request) {
	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	url := baseURL + "products"
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	// fmt.Fprintln(w, res.StatusCode)
	// fmt.Fprintln(w, string(data))
	//fmt.Println(string(data))

	// var products []Product
	// err = json.Unmarshal(data, &products)
	// tpl.ExecuteTemplate(w, "prodMain.gohtml", products)

	var products []Product
	err = json.Unmarshal(data, &products)

	if err != nil {
		log.Fatalln(err)
	}

	productsByPrice = make([]Product, 0, len(products))
	productsByQuantity = make([]Product, 0, len(products))
	productsByPrice = mergeSort(products, "Price")
	productsByQuantity = mergeSort(products, "Quantity")

	switch sortKey {
	case "":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", products)
	case "Price":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", productsByPrice)

	case "Quantity":
		tpl.ExecuteTemplate(w, "prodMain.gohtml", productsByQuantity)

	default:
		sortedProducts := sortProducts(products, sortKey)
		fmt.Println(sortedProducts)

		tpl.ExecuteTemplate(w, "prodMain.gohtml", sortedProducts)
	}
}

func sortProducts(products []Product, sortKey string) []Product {
	return mergeSort(products, sortKey)
	//return selectionSort(products, sortKey)
}

func prodAdd(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
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

		url := baseURL + "product"
		fmt.Println(url)
		_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		//direct user back to the main products page
		http.Redirect(w, r, "/products/all", http.StatusSeeOther)
		return
	}

	//get the categories & brands
	catsAndBrands := struct {
		Categories []Category
		Brands     []Brand
	}{
		Categories: clientGetCategories(),
		Brands:     clientGetBrands(),
	}
	fmt.Println(catsAndBrands)
	tpl.ExecuteTemplate(w, "prodAdd.gohtml", catsAndBrands)
}

func prodUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["productid"]

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

		url := baseURL + "product/" + productID

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			return
		}

		//direct user back to the main products page
		http.Redirect(w, r, "/products/all", http.StatusSeeOther)
	}
}

func prodDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["productid"]

	url := baseURL + "product/" + id
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	// fmt.Fprintln(w, res.StatusCode)
	// fmt.Fprintln(w, string(data))
	//fmt.Println(string(data))

	// var products []Product
	// err = json.Unmarshal(data, &products)
	// tpl.ExecuteTemplate(w, "prodMain.gohtml", products)

	var product Product
	err = json.Unmarshal(data, &product)

	if err != nil {
		log.Fatalln(err)
	}

	templateData := struct {
		Categories []Category
		Brands     []Brand
		Product    Product
	}{
		Categories: clientGetCategories(),
		Brands:     clientGetBrands(),
		Product:    product,
	}

	tpl.ExecuteTemplate(w, "prodDetail.gohtml", templateData)
}

func prodDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["productid"]

	url := baseURL + "product/" + productID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != 200 {
		return
	}

	//direct user back to the main products page
	http.Redirect(w, r, "/products/all", http.StatusSeeOther)
}

func clientGetCategories() (categories []Category) {

	url := baseURL + "categories"
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	// fmt.Fprintln(w, res.StatusCode)
	// fmt.Fprintln(w, string(data))
	fmt.Println(string(data))

	err = json.Unmarshal(data, &categories)
	return
}

func catMain(w http.ResponseWriter, r *http.Request) {
	categories := clientGetCategories()

	tpl.ExecuteTemplate(w, "catMain.gohtml", categories)
}

func catDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	catID := params["categoryid"]

	url := baseURL + "category/" + catID
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)

	var category Category
	err = json.Unmarshal(data, &category)

	if err != nil {
		log.Fatalln(err)
	}

	tpl.ExecuteTemplate(w, "catDetail.gohtml", category)
}

func catAdd(w http.ResponseWriter, r *http.Request) {

	fmt.Println("HELP HELP")
	if r.Method == http.MethodPost {

		fmt.Println("HELP")
		newCategory := Category{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(newCategory)

		url := baseURL + "category"

		_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "catAdd.gohtml", nil)
}

func catUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	catID := params["categoryid"]

	if r.Method == http.MethodPost {
		newCategory := Category{
			ID:          0,
			Name:        r.FormValue("Name"),
			Description: r.FormValue("Description"),
		}

		jsonValue, err := json.Marshal(newCategory)

		url := baseURL + "category/" + catID

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			return
		}

		//direct user back to the main products page
		http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
	}
}

func catDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	catID := params["categoryid"]

	url := baseURL + "category/" + catID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != 200 {
		return
	}

	//direct user back to the main products page
	http.Redirect(w, r, "/categories/all", http.StatusSeeOther)
}

func clientGetBrands() (brands []Brand) {

	url := baseURL + "brands"
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	// fmt.Fprintln(w, res.StatusCode)
	// fmt.Fprintln(w, string(data))
	fmt.Println("From clientGetBrands:", string(data))
	fmt.Println("---- end of data ----")

	err = json.Unmarshal(data, &brands)
	return
}
