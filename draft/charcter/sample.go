package main

import (
	"fmt"
	"unicode"

	"github.com/ikawaha/encoding/jisx0208"
)

func main() {
	for _, c := range "１二③Ⅳ" {
		fmt.Printf("%#U %v JIS X 0208 character\n", c, func() string {
			if unicode.Is(jisx0208.RangeTable, c) {
				return "is"
			}
			return "is not"
		}())
	}
	// Outpu:
	// U+FF11 '１' is a JIS X 0208 character
	// U+4E8C '二' is a JIS X 0208 character
	// U+2462 '③' is not a JIS X 0208 character
	// U+2163 'Ⅳ' is not a JIS X 0208 character
}
