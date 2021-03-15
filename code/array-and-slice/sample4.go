// +build run

package main

import "fmt"

func main() {
	slc1 := []byte{0, 1, 2, 3, 4}
	fmt.Printf("Type: %[1]T , Value: %[1]v\n", slc1)
	// Output:
	// Type: []uint8 , Value: [0 1 2 3 4]
}
