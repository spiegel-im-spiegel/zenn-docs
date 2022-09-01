//go:build run
// +build run

package main

import (
	"fmt"
	"period-sample/period"
	"time"
)

func main() {
	dt := time.Now()
	pd1 := period.New(dt)
	pd2 := period.New(dt.AddDate(0, 1, 11))
	dur := pd1.Duration(pd2)
	fmt.Println(pd2, "-", pd1, "=", dur)
	fmt.Println(pd2, "+", int(dur), "=", pd2.Add(dur))
	fmt.Println(pd1, "-", int(dur), "=", pd1.Add(-dur))
	// Output:
	// 2022年10月中旬 - 2022年9月上旬 = 4
	// 2022年10月中旬 + 4 = 2022年11月下旬
	// 2022年9月上旬 - 4 = 2022年7月下旬
}
