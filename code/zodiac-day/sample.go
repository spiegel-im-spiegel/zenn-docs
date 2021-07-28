// +build run

package main

import (
	"fmt"
	"time"
)

var (
	jst         = time.FixedZone("Asia/Tokyo", int((9 * time.Hour).Seconds()))
	zodiacNames = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	baseDay     = time.Date(2001, time.January, 1, 0, 0, 0, 0, jst)
)

func zodiacName(t time.Time) string {
	d := int64(t.Sub(baseDay).Seconds()) / 86400 % 12
	if d < 0 {
		d += 12
	}
	return zodiacNames[d]
}

func main() {
	day := time.Date(2021, time.July, 18, 0, 0, 0, 0, jst)
	for i := 0; i < 19; i++ {
		day = day.Add(time.Hour * 24)
		fmt.Printf("%v is %v\n", day.Format("2006-01-02"), zodiacName(day))
	}
}
