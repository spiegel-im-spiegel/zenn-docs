package main

import "fmt"

func main() {
	v := 1
	switch v {
	case 1:
		say := "yes"
		fmt.Println("say", say)
	case 2:
		say := "no"
		fmt.Println("say", say)
	default:
		say := "???"
		fmt.Println("say", say)
	}
}
