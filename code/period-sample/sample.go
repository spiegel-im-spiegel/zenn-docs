package main

import (
	"fmt"
	"period-sample/period"
	"time"
)

func main() {
	d := time.Now()
	pd1 := period.New(d)
	pd2 := period.New(d.AddDate(0, 1, 11))
	fmt.Println(d)
	fmt.Println(pd1)
	fmt.Println(pd2)
	fmt.Println(pd1.Add(pd2.Duration(pd1)))
	fmt.Println(pd2.Add(pd1.Duration(pd2)))
}
