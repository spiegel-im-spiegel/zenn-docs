// +build run

package main

import "fmt"

func main() {
	slc1 := []int{0, 1, 2, 3, 4}
	slc2 := slc1
	slc3 := make([]int, len(slc1), cap(slc1))
	copy(slc3, slc1)
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc2, &slc2[0], slc2)
	fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc3, &slc3[0], slc3)
	// Output:
	// Pointer: 0xc000004078 , Refer: 0xc00000c2a0 , Value: [0 1 2 3 4]
	// Pointer: 0xc000004090 , Refer: 0xc00000c2a0 , Value: [0 1 2 3 4]
	// Pointer: 0xc0000040a8 , Refer: 0xc00000c2d0 , Value: [0 1 2 3 4]
}
