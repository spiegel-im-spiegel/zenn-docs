package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	f := pflag.BoolP("foo", "f", false, "option foo")
	b := pflag.BoolP("bar", "b", false, "option bar")
	pflag.Parse()

	fmt.Println("foo = ", *f)
	fmt.Println("bar = ", *b)
}
