// +build run

package main

import "fmt"

func displayArray4Int(ary [4]int) {
	fmt.Printf("Pointer: %p , Value: %v\n", &ary, ary)
}

func main() {
	ary1 := [4]int{1, 2, 3, 4}
	ary2 := ary1

	fmt.Printf("Pointer: %p , Value: %v\n", &ary1, ary1)
	fmt.Printf("Pointer: %p , Value: %v\n", &ary2, ary2)
	displayArray4Int(ary1)
	// Output:
	// Pointer: 0xc0000141a0 , Value: [1 2 3 4]
	// Pointer: 0xc0000141c0 , Value: [1 2 3 4]
	// Pointer: 0xc000014240 , Value: [1 2 3 4]
}
