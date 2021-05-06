package main

import (
	"reflect"
)

func mergeSort(arr []Product, sortKey string) []Product {
	divPoint := len(arr) / 2
	front := arr[0:divPoint]
	back := arr[divPoint:]
	if len(front) > 1 {
		front = mergeSort(front, sortKey)
	}
	if len(back) > 1 {
		back = mergeSort(back, sortKey)
	}

	arr = merge(front, back, sortKey)

	return arr
}

func merge(front, back []Product, sortKey string) []Product {
	merged := []Product{}
	j := 0
	for i := 0; i < len(front) || j < len(back); {

		//comparing front[i] with back[j]
		var frontVal, backVal, frontField, backField reflect.Value
		if i < len(front) {
			frontVal = reflect.ValueOf(front[i])
			frontField = reflect.Indirect(frontVal).FieldByName(sortKey)
		}
		if j < len(back) {
			backVal = reflect.ValueOf(back[j])
			backField = reflect.Indirect(backVal).FieldByName(sortKey)
		}

		switch sortKey {
		case "Price":
			if i < len(front) && j < len(back) {
				if frontField.Float() <= backField.Float() {
					merged = append(merged, front[i])
					i++
				} else { // backField.Float()<frontField.Float()
					merged = append(merged, back[j])
					j++
				}

			} else if i >= len(front) && j < len(back) {
				for j < len(back) {
					merged = append(merged, back[j])
					j++
				}
			} else if i < len(front) && j >= len(back) {
				for i < len(front) {
					merged = append(merged, front[i])
					i++
				}
			}

		case "Quantity", "QuantitySold":

			if i < len(front) && j < len(back) {
				if int(frontField.Int()) <= int(backField.Int()) {
					merged = append(merged, front[i])
					i++
				} else { // int(backField.Int())<int(frontField.Int())
					merged = append(merged, back[j])
					j++
				}

			} else if i >= len(front) && j < len(back) {
				for j < len(back) {
					merged = append(merged, back[j])
					j++
				}
			} else if i < len(front) && j >= len(back) {
				for i < len(front) {
					merged = append(merged, front[i])
					i++
				}
			}

		case "Name", "DateCreated", "DateModified":

			if i < len(front) && j < len(back) {
				if frontField.String() <= backField.String() {
					merged = append(merged, front[i])
					i++
				} else { // backField.String()<frontField.String()
					merged = append(merged, back[j])
					j++
				}

			} else if i >= len(front) && j < len(back) {
				for j < len(back) {
					merged = append(merged, back[j])
					j++
				}
			} else if i < len(front) && j >= len(back) {
				for i < len(front) {
					merged = append(merged, front[i])
					i++
				}
			}

		} //end switch sortKey

	} //end for
	return merged
}
