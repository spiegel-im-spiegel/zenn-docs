package period

import (
	"fmt"
	"time"
)

// Period is period count in a month
type Period int

const (
	Early Period = iota + 1
	Middle
	Late
	maxPeriods    = int(Late)
	periodsByYear = 12 * maxPeriods
)

var periodmap = map[Period]string{
	Early:  "上旬",
	Middle: "中旬",
	Late:   "下旬",
}

func (p Period) String() string {
	if s, ok := periodmap[p]; ok {
		return s
	}
	return ""
}

// Date is date value with Period.
type Date int

// Duration is duration value from Date to date.
type Duration int

// New creates Date instance from time.Time.
func New(dt time.Time) Date {
	d := dt.Year()*periodsByYear + (int(dt.Month())-1)*maxPeriods
	if dt.Day() > 10 {
		d++
	}
	if dt.Day() > 20 {
		d++
	}
	return Date(d)
}

// Year returns year from Date.
func (pd Date) Year() int {
	return int(pd) / periodsByYear
}

// Month returns month from Date.
func (pd Date) Month() int {
	return (int(pd)%periodsByYear)/maxPeriods + 1
}

// Year returns Period value from Date.
func (pd Date) Period() Period {
	return Period(((int(pd) % periodsByYear) % maxPeriods) + 1)
}

func (pd Date) String() string {
	return fmt.Sprintf("%d年%d月%v", pd.Year(), pd.Month(), pd.Period())
}

// Duration returns Duration value from start to end.
func (start Date) Duration(end Date) Duration {
	return Duration(int(end) - int(start))
}

// Add returns result value of add operation.
func (pd Date) Add(dur Duration) Date {
	return Date(int(pd) + int(dur))
}
