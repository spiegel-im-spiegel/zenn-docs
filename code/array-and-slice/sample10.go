// +build run

package main

import (
	"fmt"
	"reflect"
)

func main() {
	slc1 := []int{0, 1, 2, 3, 4}
	slc2 := []int{0, 1, 2, 3, 4}
	// fmt.Printf("slc1 == slc2: %v\n", slc1 == slc2) // invalid operation: slc1 == slc2 (slice can only be compared to nil)
	if reflect.DeepEqual(slc1, slc2) {
		fmt.Println("slc1 == slc2: true")
	} else {
		fmt.Println("slc1 == slc2: false")
	}
	// Output
	// slc1 == slc2: true
}
