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

type ErrorTplData struct {
	ErrorMessage string
	RedirectLink string
	ButtonValue  string
}

var (
	storedProducts   []Product
	storedCategories []Category
	storedBrands     []Brand

	productsByPrice    []Product
	productsByQuantity []Product
	// productsByName     []Product

	shopMap = make(map[string]shopCart)

	results []Product
)

const sessionID string = "abcde"

func initGlobalVars() {
	storedProducts, _ = clientGetProducts()
	storedCategories, _ = clientGetCategories()
	storedBrands, _ = clientGetBrands()

	fmt.Println("Products:", storedProducts)
	fmt.Println()
	fmt.Println("Categories:", storedCategories)
	fmt.Println()
	fmt.Println("Brands:", storedBrands)

	productsByPrice = make([]Product, 0, len(storedProducts))
	productsByQuantity = make([]Product, 0, len(storedProducts))
	productsByPrice = mergeSort(storedProducts, "Price")
	productsByQuantity = mergeSort(storedProducts, "Quantity")

	// productsByName = make([]Product, 0, len(storedProducts))
	// productsByName = mergeSort(storedProducts, "Name")
}

// const baseURL = "http://localhost:5000/api/v1/admin/"

func UserSearch(w http.ResponseWriter, r *http.Request) {
	//check if global vars are initialized.
	if len(storedProducts) == 0 {
		initGlobalVars()
	}

	type templateData struct {
		Products   []prodTpl
		Categories []Category
		Brands     []Brand
	}

	// fmt.Println("storedBrands=", storedBrands)
	// fmt.Println("storedCategories=", storedCategories)
	// fmt.Println("storedProducts=", storedProducts)

	sortKey := r.FormValue("sortby")
	fmt.Println("sortKey =", sortKey)

	//search the global var storedProducts for active products
	var activeProducts []Product
	for _, v := range storedProducts {
		if v.Status == "active" {
			activeProducts = append(activeProducts, v)
		}
	}
	fmt.Println("From UserSearch: activeProducts = ", activeProducts)

	if r.Method == http.MethodPost {
		searchKey := r.FormValue("SearchKey")
		catID, _ := strconv.Atoi(r.FormValue("CatID"))
		brandID, _ := strconv.Atoi(r.FormValue("BrandID"))
		results = searchProduct(searchKey, catID, brandID, activeProducts)
	}

	switch sortKey {
	case "":
		tpl.ExecuteTemplate(w, "userSearch.gohtml", templateData{formatProdTplData(results), storedCategories, storedBrands})
	default:
		sortedResults := sortProducts(results, sortKey)
		fmt.Println(sortedResults)

		tpl.ExecuteTemplate(w, "userSearch.gohtml", templateData{formatProdTplData(sortedResults), storedCategories, storedBrands})
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	initGlobalVars()

	type templateData struct {
		Products   []prodTpl
		Categories []Category
		Brands     []Brand
	}

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

	if r.Method == http.MethodPost && r.FormValue("Search") == "Search" {
		searchKey := r.FormValue("SearchKey")
		catID, _ := strconv.Atoi(r.FormValue("CatID"))
		brandID, _ := strconv.Atoi(r.FormValue("BrandID"))
		results := searchProduct(searchKey, catID, brandID, activeProducts)
		tpl.ExecuteTemplate(w, "userSearch.gohtml", templateData{formatProdTplData(results), storedCategories, storedBrands})
		return
	}

	switch sortKey {
	case "":
		tpl.ExecuteTemplate(w, "index.gohtml", templateData{formatProdTplData(activeProducts), storedCategories, storedBrands})
	case "Price":
		//search the global var productsByPrice for active products
		var activeProductsByPrice []Product
		for _, v := range productsByPrice {
			if v.Status == "active" {
				activeProductsByPrice = append(activeProductsByPrice, v)
			}
		}
		tpl.ExecuteTemplate(w, "index.gohtml", templateData{formatProdTplData(activeProductsByPrice), storedCategories, storedBrands})

	case "Quantity":
		//search the global var productsByQuantity for active products
		var activeProductsByQuantity []Product
		for _, v := range productsByQuantity {
			if v.Status == "active" {
				activeProductsByQuantity = append(activeProductsByQuantity, v)
			}
		}
		tpl.ExecuteTemplate(w, "index.gohtml", templateData{formatProdTplData(activeProductsByQuantity), storedCategories, storedBrands})

	default:
		sortedProducts := sortProducts(activeProducts, sortKey)
		fmt.Println(sortedProducts)

		tpl.ExecuteTemplate(w, "index.gohtml", templateData{formatProdTplData(sortedProducts), storedCategories, storedBrands})
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

		for i, v := range userCart {
			if v.ID == id {
				addedToCart = true

				if v.QuantityToBuy+qty > product.Quantity {
					addedToCartMsg = "Quantity Exceeded! Failed to add to cart."
				} else {
					v.QuantityToBuy += qty
					fmt.Println("Adding same product: ", userCart[i])
					userCart[i] = v
					fmt.Println("Adding same product - after assigning v to userCart[i]: ", userCart[i])

					addedToCartMsg = fmt.Sprintf("%d are added to your cart.", qty)
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

	if r.Method == http.MethodPost && r.FormValue("submit") == "Remove From Cart" {
		cartIndex, _ := strconv.Atoi(r.FormValue("cartIndex"))

		// copy(userCart[cartIndex:], userCart[cartIndex+1:])
		// userCart[len(userCart)-1] = CartItem{}
		// userCart = userCart[:len(userCart)-1]

		userCart = append(userCart[:cartIndex], userCart[cartIndex+1:]...)
		shopMap[sessionID] = userCart

	} else if r.Method == http.MethodPost && r.FormValue("submit") == "Checkout" {
		fmt.Println("Checking out.........")
		type cartItemWithPrice struct {
			ID            int
			Name          string
			Price         float64
			QuantityToBuy int
			ItemTotal     float64
		}
		type templateStruct struct {
			CartData    []cartItemWithPrice
			CartTotal   float64
			CartIsEmpty bool
		}

		var emptyCart bool
		if len(userCart) == 0 {
			emptyCart = true
			//load template
			tplData := templateStruct{[]cartItemWithPrice{}, 0, emptyCart}
			tpl.ExecuteTemplate(w, "orderConfirmation.gohtml", tplData)
			return
		}

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

		res, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			return
		}

		//print out the invoice
		cartWithPrice := []cartItemWithPrice{}
		var cartTotal float64
		for _, v := range userCart {
			itemTotal := v.Price * float64(v.QuantityToBuy)
			cartTotal += itemTotal
			cartWithPrice = append(cartWithPrice, cartItemWithPrice{v.ID, v.Name, v.Price, v.QuantityToBuy, itemTotal})
		}
		tplData := templateStruct{cartWithPrice, cartTotal, emptyCart}

		//empty the cart
		var emptiedCart shopCart
		shopMap[sessionID] = emptiedCart

		// execute order confirmation page
		tpl.ExecuteTemplate(w, "orderConfirmation.gohtml", tplData)
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

	tpl.ExecuteTemplate(w, "cart.gohtml", templateData)
}

func Enquiry(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("Name")
		email := r.FormValue("Email")
		message := r.FormValue("Message")

		enquiry := struct {
			Name        string
			Email       string
			EnquiryDate string
			Message     string
		}{
			Name:        name,
			Email:       email,
			EnquiryDate: "",
			Message:     message,
		}
		jsonValue, err := json.Marshal(enquiry)
		if err != nil {
			log.Fatalln(err)
		}

		url := baseURL + "enquiry"
		res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			return
		}
		adminEmail := "yenyenjoe@yahoo.com.sg"
		msg := []byte("To: " + adminEmail + "\r\n" +
			"Subject: Enquiry on GoGreen from " + enquiry.Name + "\r\n" +
			"\r\n" +
			enquiry.Message + "\r\n Contact email: " + enquiry.Email +
			"\r\n")

		sendMail(msg)

		// execute order confirmation page
		tpl.ExecuteTemplate(w, "enquiryConfirmation.gohtml", nil)
		return
	}

	// execute order confirmation page
	tpl.ExecuteTemplate(w, "enquiry.gohtml", nil)
}
