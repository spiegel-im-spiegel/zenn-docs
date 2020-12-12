package main

import "fmt"

type Number struct {
	n int
}

func (n Number) Add(i int) int {
	return n.n + i
}

func main() {
	increment := Number{1}.Add
	fmt.Printf("%T\n", increment)
	fmt.Printf("%#x\n", 32<<(^uint32(0)>>63))
	fmt.Println(increment(2))
}
