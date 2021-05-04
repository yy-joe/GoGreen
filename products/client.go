package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	sortedProducts := sortProducts(products, sortKey)
	fmt.Println("-----------------sortedProducts :----------------------")
	fmt.Println(sortedProducts)
	tpl.ExecuteTemplate(w, "prodMain.gohtml", sortedProducts)
}

func sortProducts(products []Product, sortKey string) []Product { //currently using selection sort

	return selectionSort(products, sortKey)
}

// type Product struct {
// 	ID           int
// 	Name         string
// 	Image        string
// 	DescShort    string
// 	DescLong     string
// 	DateCreated  string
// 	DateModified string
// 	Price        float64
// 	Quantity     int
// 	Condition    string
// 	CategoryID   int
// 	BrandID      int
// 	Status       string
// }

func prodDetail(w http.ResponseWriter, r *http.Request) {

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
	tpl.ExecuteTemplate(w, "prodDetail.gohtml", catsAndBrands)
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
