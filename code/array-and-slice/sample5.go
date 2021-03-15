// +build run

package main

import "fmt"

func main() {
	ary1 := [5]byte{0, 1, 2, 3, 4}
	slc1 := ary1[:]
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &ary1, &ary1[0], ary1)
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
	// Output:
	// Pointer: 0xc000012088 , Refer: 0xc000012088 , Value: [0 1 2 3 4]
	// Pointer: 0xc000004078 , Refer: 0xc000012088 , Value: [0 1 2 3 4]
}
