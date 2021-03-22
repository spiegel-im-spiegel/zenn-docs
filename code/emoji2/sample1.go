// +build run

package main

import (
	"fmt"

	"github.com/kyokomi/emoji/v2"
)

func main() {
	fmt.Println("Hello World Emoji!")

	emoji.Println(":beer: Beer!!!")

	pizzaMessage := emoji.Sprint("I like a :pizza: and :sushi:!!")
	fmt.Println(pizzaMessage)
}
