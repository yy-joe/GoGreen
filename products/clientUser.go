package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CartItem struct {
	ID            int
	Name          string
	Price         float64
	QuantityToBuy int
}

var (
	storedProducts   []Product
	storedCategories []Category
	storedBrands     []Brand

	productsByPrice    []Product
	productsByQuantity []Product
	productsByName     []Product

	shoppingCart []CartItem
)

func initGlobalVars() {
	storedProducts = clientGetProducts()
	storedCategories = clientGetCategories()
	storedBrands = clientGetBrands()

	fmt.Println("Products:", storedProducts)
	fmt.Println()
	fmt.Println("Categories:", storedCategories)
	fmt.Println()
	fmt.Println("Brands:", storedBrands)

	productsByPrice = make([]Product, 0, len(storedProducts))
	productsByQuantity = make([]Product, 0, len(storedProducts))
	productsByPrice = mergeSort(storedProducts, "Price")
	productsByQuantity = mergeSort(storedProducts, "Quantity")

	productsByName = make([]Product, 0, len(storedProducts))
	productsByName = mergeSort(storedProducts, "Name")
}

// const baseURL = "http://localhost:5000/api/v1/admin/"

func index(w http.ResponseWriter, r *http.Request) {
	initGlobalVars()

	fmt.Println("storedBrands=", storedBrands)
	fmt.Println("storedCategories=", storedCategories)
	fmt.Println("storedProducts=", storedProducts)

	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	//search the global var storedProducts for active products
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
	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}
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

	var addedToCart bool
	var addedToCartMsg string
	if r.Method == http.MethodPost {
		id := r.FormValue("productid")
		price, _ := strconv.ParseFloat(r.FormValue("productprice"), 64)
		qty, _ := strconv.Atoi(r.FormValue("quantityToBuy"))
		fmt.Println(id, price)

		//update shopping cart with these values

		//pass added to cart message to the template
		addedToCart = true
		addedToCartMsg = fmt.Sprintf("%d are added to your cart.", qty)
	}

	templateData := struct {
		CategoryName   string
		BrandName      string
		Product        Product
		AddedToCart    bool
		AddedToCartMsg string
	}{
		CategoryName:   categoryName,
		BrandName:      brandName,
		Product:        product,
		AddedToCart:    addedToCart,
		AddedToCartMsg: addedToCartMsg,
	}

	tpl.ExecuteTemplate(w, "details.gohtml", templateData)
}
