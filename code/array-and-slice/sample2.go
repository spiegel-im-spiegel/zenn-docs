// +build run

package main

import "fmt"

func main() {
	ary1 := [4]int{1, 2, 3, 4}
	ary2 := [4]int{1, 2, 3, 4}
	ary3 := [4]int{2, 3, 4, 5}
	ary4 := [4]int64{1, 2, 3, 4}

	fmt.Printf("ary1 == ary2: %v\n", ary1 == ary2) // ary1 == ary2: true
	fmt.Printf("ary1 == ary3: %v\n", ary1 == ary3) // ary1 == ary3: false
	fmt.Printf("ary1 == ary4: %v\n", ary1 == ary4) // invalid operation: ary1 == ary4 (mismatched types [4]int and [4]int64)
}
