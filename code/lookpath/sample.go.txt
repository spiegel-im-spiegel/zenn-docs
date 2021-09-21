package main

import (
	"os"
	"fmt"

	"golang.org/x/sys/execabs"
)

func main() {
	if b, err := execabs.Command("hello").Output(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println("Say:", string(b))
	}
}
