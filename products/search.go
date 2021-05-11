package products

import "strings"

// var (
// 	storedProducts   []Product
// 	storedCategories []Category
// 	storedBrands     []Brand

// 	productsByPrice    []Product
// 	productsByQuantity []Product
// 	productsByName     []Product

// 	shopMap = make(map[string]shopCart)
// )
func searchProduct(searchKey string, catID int, brandID int, products []Product) (results []Product) {
	// fmt.Println(strings.Contains("seafood restaurant", "foo"))
	for _, v := range products {
		if strings.Contains(v.Name, searchKey) {
			if catID <= 0 && brandID <= 0 {
				results = append(results, v)
			} else if catID <= 0 && brandID == v.BrandID {
				results = append(results, v)
			} else if catID == v.CategoryID && brandID <= 0 {
				results = append(results, v)
			} else if catID == v.CategoryID && brandID == v.BrandID {
				results = append(results, v)
			}
		}
	}
	return
}
