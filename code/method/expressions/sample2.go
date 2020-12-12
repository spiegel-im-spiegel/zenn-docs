package main

import "fmt"

type Number struct {
	n int
}

func (n Number) Add(i int) int {
	return n.n + i
}

func (n Number) Sub(i int) int {
	return n.n - i
}

func main() {
	add := Number.Add
	fmt.Printf("%T\n", add)
	sub := Number.Sub
	fmt.Printf("%T\n", sub)
	fmt.Println(add(Number{1}, 2)) //Output: 3
	fmt.Println(sub(Number{1}, 2)) //Output: -1
}
