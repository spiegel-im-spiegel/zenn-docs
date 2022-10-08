package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

func main() {
	f := flag.Bool("foo", false, "option foo")
	b := flag.Bool("bar", false, "option bar")
	flag.Parse()

	fmt.Println("foo = ", *f)
	fmt.Println("bar = ", *b)
}
