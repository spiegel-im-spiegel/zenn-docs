package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/dimchansky/utfbom"
)

const text = "\xEF\xBB\xBFhello"

func main() {
	fmt.Println([]byte(text))
	r := utfbom.SkipOnly(strings.NewReader(text))
	b, err := io.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
}
