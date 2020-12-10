package main

import (
	"errors"
	"fmt"
	)

func main() {
	err := bar()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	fmt.Println("Normal End.")
}

func bar() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("Recovered from: %w", rec)
		}
	}()

	foo()
	return
}

func foo() {
	panic(errors.New("Panic!"))
}
