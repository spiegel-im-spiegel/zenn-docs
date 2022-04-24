package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goark/csvdata"
	"golang.org/x/text/encoding/japanese"
)

func main() {
	rc := csvdata.NewRows(csvdata.New(japanese.ShiftJIS.NewDecoder().Reader(os.Stdin)), true)
	for {
		if err := rc.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Fprintln(os.Stderr, err)
			return
		}
		order, err := rc.ColumnInt64("order", 10)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("    Order =", order)
		fmt.Println("     Name =", rc.Column("name"))
		mass, err := rc.ColumnFloat64("mass")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("     Mass =", mass)
		habitable, err := rc.ColumnBool("habitable")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Habitable =", habitable)
	}
}
