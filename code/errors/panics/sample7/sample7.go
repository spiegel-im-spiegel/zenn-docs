package main

import "fmt"

func main() {
	foo()
}

func foo() {
	numbers := []int{0, 1, 2}
	fmt.Println(numbers[3])
}
