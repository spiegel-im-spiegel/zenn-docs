// +build run

package main

import "fmt"

func displaySliceByte(slc []byte) {
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc, &slc[0], slc)
}

func main() {
	ary1 := [5]byte{0, 1, 2, 3, 4}
	slc1 := ary1[:]
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &ary1, &ary1[0], ary1)
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
	displaySliceByte(slc1)
	// Output:
	// Pointer: 0xc000102058 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
	// Pointer: 0xc000100048 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
	// Pointer: 0xc000100078 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
}
