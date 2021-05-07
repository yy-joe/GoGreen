package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// const baseURL = "http://localhost:5000/api/v1/admin/"

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("storedBrands=", storedBrands)
	fmt.Println("storedCategories=", storedCategories)
	fmt.Println("storedProducts=", storedProducts)

	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	//search the global var storedProducts for active products
	fmt.Println("a=", a)
	var activeProducts []Product
	for _, v := range storedProducts {
		if v.Status == "active" {
			activeProducts = append(activeProducts, v)
		}
	}
	fmt.Println(activeProducts)

	//search the global var productsByPrice for active products
	var activeProductsByPrice []Product
	for _, v := range productsByPrice {
		if v.Status == "active" {
			activeProductsByPrice = append(activeProductsByPrice, v)
		}
	}

	// //search the global var productsByQuantity for active products
	// var activeProductsByQuantity []Product
	// for _, v := range productsByQuantity {
	// 	if v.Status == "active" {
	// 		activeProductsByQuantity = append(activeProductsByQuantity, v)
	// 	}
	// }

	switch sortKey {
	case "":
		tpl.ExecuteTemplate(w, "index.gohtml", activeProducts)
	case "Price":
		tpl.ExecuteTemplate(w, "index.gohtml", activeProductsByPrice)

	// case "Quantity":
	// 	tpl.ExecuteTemplate(w, "index.gohtml", activeProductsByQuantity)

	default:
		sortedProducts := sortProducts(activeProducts, sortKey)
		fmt.Println(sortedProducts)

		tpl.ExecuteTemplate(w, "index.gohtml", sortedProducts)
	}
}

func details(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["productid"])

	var product Product
	for _, v := range storedProducts {
		if v.ID == id {
			product = v
		}
	}

	//get the category name
	var categoryName string
	for _, v := range storedCategories {
		if v.ID == product.CategoryID {
			categoryName = v.Name
		}
	}

	//get the brand name
	var brandName string
	for _, v := range storedBrands {
		if v.ID == product.BrandID {
			brandName = v.Name
		}
	}

	templateData := struct {
		CategoryName string
		BrandName    string
		Product      Product
	}{
		CategoryName: categoryName,
		BrandName:    brandName,
		Product:      product,
	}

	tpl.ExecuteTemplate(w, "details.gohtml", templateData)
}
