package main

import (
	"fmt"
	"sort"
)

func main() {
	ds := []float64{0.055, 0.815, 1.0, 0.107}
	fmt.Println(ds) //before
	sort.Slice(ds, func(i, j int) bool {
		return ds[i] < ds[j]
	})
	fmt.Println(ds) //after
}
