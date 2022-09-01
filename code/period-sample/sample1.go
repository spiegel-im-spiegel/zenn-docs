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
	fmt.Println(dt.Format("2006-01-02"), "->", period.New(dt))
	// Output:
	// 2022-09-01 -> 2022年9月上旬
}
