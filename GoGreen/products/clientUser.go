package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
type shopCart []CartItem

var (
	storedProducts   []Product
	storedCategories []Category
	storedBrands     []Brand

	productsByPrice    []Product
	productsByQuantity []Product
	productsByName     []Product

	shopMap = make(map[string]shopCart)
)

const sessionID string = "abcde"

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

func Index(w http.ResponseWriter, r *http.Request) {
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

func Details(w http.ResponseWriter, r *http.Request) {
	//get the userID/ sessionID

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
	var userCart shopCart
	if r.Method == http.MethodPost {
		id, _ := strconv.Atoi(r.FormValue("productid"))
		name := r.FormValue("productname")
		price, _ := strconv.ParseFloat(r.FormValue("productprice"), 64)
		qty, _ := strconv.Atoi(r.FormValue("quantityToBuy"))
		fmt.Println(id, price)

		//update shopping cart with these values
		userCart = shopMap[sessionID]

		for _, v := range userCart {
			if v.ID == id {
				addedToCart = true

				if v.QuantityToBuy+qty > product.Quantity {
					addedToCartMsg = "Quantity Exceeded! Failed to add to cart."
				} else {
					v.QuantityToBuy += qty
				}
			}
		}

		if !addedToCart {
			userCart = append(userCart, CartItem{id, name, price, qty})
			addedToCart = true
			addedToCartMsg = fmt.Sprintf("%d are added to your cart.", qty)
		}

		shopMap[sessionID] = userCart

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

func Cart(w http.ResponseWriter, r *http.Request) {
	//get the userID/ sessionID

	userCart := shopMap[sessionID]

	if r.Method == http.MethodPost {

		jsonValue, err := json.Marshal(userCart)

		if err != nil {
			log.Fatalln(err)
		}

		url := baseURL + "product/quantity-update"

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

		// execute order confirmation page
		tpl.ExecuteTemplate(w, "orderConfirmation.gohtml", nil)
		return
	}

	type cartItemWithPrice struct {
		ID            int
		Name          string
		Price         float64
		QuantityToBuy int
		ItemTotal     float64
	}
	cartWithPrice := []cartItemWithPrice{}
	var cartTotal float64
	for _, v := range userCart {
		itemTotal := v.Price * float64(v.QuantityToBuy)
		cartTotal += itemTotal
		cartWithPrice = append(cartWithPrice, cartItemWithPrice{v.ID, v.Name, v.Price, v.QuantityToBuy, itemTotal})
	}
	templateData := struct {
		CartData  []cartItemWithPrice
		CartTotal float64
	}{
		CartData:  cartWithPrice,
		CartTotal: cartTotal,
	}

	fmt.Println("============================================================")
	fmt.Println(shopMap[sessionID])
	fmt.Println("============================================================")

	tpl.ExecuteTemplate(w, "cart.gohtml", templateData)
}
