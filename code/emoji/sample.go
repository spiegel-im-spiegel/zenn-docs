package main

import (
    "fmt"
    "unicode"
)

func check(r rune) string {
    switch {
    case unicode.Is(unicode.Variation_Selector, r):
        return "Variation Selector"
    case unicode.Is(unicode.Sc, r):
        return "Symbol/currency"
    case unicode.Is(unicode.Sk, r):
        return "Symbol/modifier"
    case unicode.Is(unicode.Sm, r):
        return "Symbol/math"
    case unicode.Is(unicode.So, r):
        return "Symbol/other"
    case unicode.Is(unicode.Lm, r):
        return "Letter/modifier"
    case unicode.Is(unicode.Lo, r):
        return "Letter/other"
    case unicode.Is(unicode.Nl, r):
        return "Number/letter"
    case unicode.Is(unicode.No, r):
        return "Number/other"
    case unicode.Is(unicode.Mc, r):
        return "Mark/spacing combining"
    case unicode.Is(unicode.Me, r):
        return "Mark/enclosing"
    case unicode.Is(unicode.Mn, r):
        return "Mark/nonspacing"
    case unicode.Is(unicode.Pc, r):
        return "Punctuation/connector"
    case unicode.Is(unicode.Pd, r):
        return "Punctuation/dash"
    case unicode.Is(unicode.Pe, r):
        return "Punctuation/close"
    case unicode.Is(unicode.Pf, r):
        return "Punctuation/final quote"
    case unicode.Is(unicode.Pi, r):
        return "Punctuation/initial quote"
    case unicode.Is(unicode.Ps, r):
        return "Punctuation/open"
    case unicode.Is(unicode.Po, r):
        return "Punctuation/other"
    case unicode.Is(unicode.Zl, r):
        return "Separator/line"
    case unicode.Is(unicode.Zp, r):
        return "Separator/paragraph"
    case unicode.Is(unicode.Zs, r):
        return "Separator/space"
    case unicode.IsGraphic(r):
        return "Graphic"
    case unicode.Is(unicode.Join_Control, r):
        return "Join Control"
    case unicode.Is(unicode.Cc, r):
        return "Control/control"
    case unicode.Is(unicode.Cf, r):
        return "Control/format"
    case unicode.Is(unicode.Cs, r):
        return "Control/surrogate"
    case unicode.Is(unicode.Co, r):
        return "Control/private use"
    }
    return "Unknown"
}

func main() {
    text := "㍻\t⌨\n🤔 #️⃣ 🙇‍♂️"
    for _, c := range text {
        fmt.Printf("%#U (%v)\n", c, check(c))
    }
}
