package products

import (
	"fmt"
	"reflect"
)

func selectionSort(arr []Product, sortKey string) []Product {
	if len(arr) == 0 {
		fmt.Println("Empty array!")
		return nil
	} else if len(arr) == 1 {
		return arr
	}

	//array length is at least 2
	last := len(arr) - 1
	for last > 0 {
		biggestIndex := 0

		for j := 1; j <= last; j++ {
			jVal := reflect.ValueOf(arr[j])
			jField := reflect.Indirect(jVal).FieldByName(sortKey)
			biggestVal := reflect.ValueOf(arr[biggestIndex])
			biggestField := reflect.Indirect(biggestVal).FieldByName(sortKey)

			switch sortKey {
			case "Price":
				if jField.Float() > biggestField.Float() {
					biggestIndex = j
				}
			case "Quantity", "Sales":
				if int(jField.Int()) > int(biggestField.Int()) {
					biggestIndex = j
				}
			case "Name", "DateCreated", "DateModified":
				if jField.String() > biggestField.String() {
					biggestIndex = j
				}
			}

		}

		if biggestIndex != last {
			temp := arr[last]
			arr[last] = arr[biggestIndex]
			arr[biggestIndex] = temp
		}
		last--
	}
	return arr
}

func selectionSort_working(arr []Product) []Product {
	if len(arr) == 0 {
		fmt.Println("Empty array!")
		return nil
	} else if len(arr) == 1 {
		return arr
	}

	//array length is at least 2
	last := len(arr) - 1
	for last > 0 {
		biggestIndex := 0

		for j := 1; j <= last; j++ {
			if (arr[j]).Price > (arr[biggestIndex]).Price {
				biggestIndex = j
			}
		}

		if biggestIndex != last {
			temp := arr[last]
			arr[last] = arr[biggestIndex]
			arr[biggestIndex] = temp
		}
		last--
	}
	return arr
}
