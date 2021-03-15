// +build run

package main

import "fmt"

func main() {
	var slc []int
	fmt.Printf("Pointer: %p , <ZERO value>\n", &slc)
	for i := 0; i < 5; i++ {
		slc = append(slc, i)
		fmt.Printf("Pointer: %p , Refer: %p , Value: %v (%d)\n", &slc, &slc[0], slc, cap(slc))
	}
	// Output:
	// Pointer: 0xc000004078 , <ZERO value>
	// Pointer: 0xc000004078 , Refer: 0xc000012088 , Value: [0] (1)
	// Pointer: 0xc000004078 , Refer: 0xc0000120d0 , Value: [0 1] (2)
	// Pointer: 0xc000004078 , Refer: 0xc0000141c0 , Value: [0 1 2] (4)
	// Pointer: 0xc000004078 , Refer: 0xc0000141c0 , Value: [0 1 2 3] (4)
	// Pointer: 0xc000004078 , Refer: 0xc00000e340 , Value: [0 1 2 3 4] (8)
}
