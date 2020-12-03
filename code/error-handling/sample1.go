package main

import (
	"fmt"
	"os"
)

func main() {
	filename := "go.mod"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
}
