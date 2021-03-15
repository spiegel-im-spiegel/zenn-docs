// +build run

package main

import "fmt"

func main() {
	ary1 := [4]int{1, 2, 3, 4}
	fmt.Printf("Type: %[1]T , Value: %[1]v\n", ary1)
	// Output:
	// Type: [4]int , Value: [1 2 3 4]
}
