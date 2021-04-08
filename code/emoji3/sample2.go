// +build run

package main

import (
	"fmt"

	"github.com/rivo/uniseg"
)

func main() {
	text := "ペンギン ﾍﾟﾝｷﾞﾝ"
	fmt.Println("Text:", text)
	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		rs := gr.Runes()
		fmt.Printf("%v : %U\n", string(rs), rs)
	}
}
