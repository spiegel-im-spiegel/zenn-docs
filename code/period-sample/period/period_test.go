package period_test

import (
	"period-sample/period"
	"testing"
	"time"
)

func TestPeriod(t *testing.T) {
	testCases := []struct {
		p   period.Period
		str string
	}{
		{p: period.Period(0), str: ""},
		{p: period.Early, str: "上旬"},
		{p: period.Middle, str: "中旬"},
		{p: period.Late, str: "下旬"},
		{p: period.Period(4), str: ""},
	}
	for _, tc := range testCases {
		if tc.p.String() != tc.str {
			t.Errorf("Period(%[1]d) = \"%[1]v\", want \"%[2]v\".", tc.p, tc.str)
		}
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		dt  time.Time
		str string
	}{
		{dt: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC), str: "2022年9月上旬"},
		{dt: time.Date(2022, time.September, 10, 0, 0, 0, 0, time.UTC), str: "2022年9月上旬"},
		{dt: time.Date(2022, time.September, 11, 0, 0, 0, 0, time.UTC), str: "2022年9月中旬"},
		{dt: time.Date(2022, time.September, 20, 0, 0, 0, 0, time.UTC), str: "2022年9月中旬"},
		{dt: time.Date(2022, time.September, 21, 0, 0, 0, 0, time.UTC), str: "2022年9月下旬"},
		{dt: time.Date(2022, time.September, 30, 0, 0, 0, 0, time.UTC), str: "2022年9月下旬"},
	}
	for _, tc := range testCases {
		dt := period.New(tc.dt)
		if dt.String() != tc.str {
			t.Errorf("New(%v) = \"%v\", want \"%v\".", tc.dt, dt, tc.str)
		}
	}
}

func TestDuration(t *testing.T) {
	testCases := []struct {
		dt1 time.Time
		dt2 time.Time
		dur period.Duration
	}{
		{dt1: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC), dt2: time.Date(2022, time.October, 11, 0, 0, 0, 0, time.UTC), dur: 4},
		{dt1: time.Date(2022, time.October, 11, 0, 0, 0, 0, time.UTC), dt2: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC), dur: -4},
	}
	for _, tc := range testCases {
		dur := period.New(tc.dt1).Duration(period.New(tc.dt2))
		if dur != tc.dur {
			t.Errorf("Duration(%v - %v) = %v, want %v.", tc.dt1, tc.dt2, dur, tc.dur)
		}
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		dt  time.Time
		dur period.Duration
		res period.Date
	}{
		{dt: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC), dur: 4, res: period.Date(72820)},
		{dt: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC), dur: -4, res: period.Date(72812)},
	}
	for _, tc := range testCases {
		res := period.New(tc.dt).Add(tc.dur)
		if res != tc.res {
			t.Errorf("Date(%v).Add(%v) = %v, want %v.", tc.dt, tc.dur, res, tc.res)
		}
	}
}
