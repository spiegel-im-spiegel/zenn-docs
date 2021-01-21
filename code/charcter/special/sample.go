package main

import (
	"fmt"
	"unicode"
)

func check(r rune) string {
	switch {
	case unicode.Is(unicode.Radical, r): return "Radical (Symbol)"
	case unicode.Is(unicode.Ideographic, r): return "Ideographic (Letter)"
	case unicode.Is(unicode.Variation_Selector, r): return "Variation Selector (Mark)"
	case unicode.Is(unicode.Join_Control, r): return "Join Control"
	}
	switch {
	case unicode.IsSymbol(r): return "Symbol"
	case unicode.IsLetter(r): return "Letter"
	case unicode.IsNumber(r): return "Number"
	case unicode.IsMark(r): return "Mark"
	case unicode.IsPunct(r): return "Punct"
	case unicode.IsSpace(r): return "Space"
	case unicode.IsControl(r): return "Control"
	}
	if unicode.IsGraphic(r) {return "Graphic"}
	return "Unknown"
}

func main() {
	for _, c := range []rune{0x0001f647, 0x200d, 0x2642, 0xfe0f, 0x30D8, 0x309A, 0xFF8D, 0xFF9F, 0x2F5F, 0x7389, 0x908A, 0xE0104, 0x3002, 0xff11, 0x0031, 0x0022, 0x0027, 0x201d, 0x300c, 0x300d} {
		fmt.Printf("%#U (%v)\n", c, check(c))
	}
}
