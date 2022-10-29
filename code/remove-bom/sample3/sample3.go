package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spkg/bom"
)

func main() {
	file, err := os.Open("./sample3.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	r := csv.NewReader(bom.NewReader(file))
	for {
		row, err := r.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			break
		}
		fmt.Println(row)
	}
}
