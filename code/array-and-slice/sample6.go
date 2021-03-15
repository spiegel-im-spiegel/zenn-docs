// +build run

package main

import "fmt"

func main() {
	ary1 := [5]byte{0, 1, 2, 3, 4}
	slc1 := ary1[:]
	slc2 := ary1[2:4]
	slc3 := slc2[:cap(slc2)]
	fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &ary1[0], len(ary1), cap(ary1), ary1)
	fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc1[0], len(slc1), cap(slc1), slc1)
	fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc2[0], len(slc2), cap(slc2), slc2)
	fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc3[0], len(slc3), cap(slc3), slc3)
	// Output:
	// Refer: 0xc000012088 , Len: 5 , Cap: 5 , Value: [0 1 2 3 4]
	// Refer: 0xc000012088 , Len: 5 , Cap: 5 , Value: [0 1 2 3 4]
	// Refer: 0xc00001208a , Len: 2 , Cap: 3 , Value: [2 3]
	// Refer: 0xc00001208a , Len: 3 , Cap: 3 , Value: [2 3 4]
}
