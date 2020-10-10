package main

import "fmt"

func main() {
	var a *string = nil
	var b interface{} = a
	fmt.Println("a == nil:", a == nil) // true
	fmt.Println("b == nil:", b == nil) // false
	fmt.Println("a == b:  ", a == b)   // true
}
