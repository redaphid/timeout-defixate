package main

import (
	"testing"
	"time"
)

func at(hour int) time.Time {
	return time.Date(2026, 6, 18, hour, 30, 0, 0, time.Local)
}

func TestInCurfewWrapsMidnight(t *testing.T) {
	cases := []struct {
		hour int
		want bool
	}{
		{21, false}, // just before start
		{22, true},  // start edge, inclusive
		{23, true},
		{0, true}, // past midnight
		{7, true},
		{8, false}, // end edge, exclusive
		{9, false},
		{12, false},
	}

	for _, c := range cases {
		got := inCurfew(at(c.hour), 22, 8)
		if got != c.want {
			t.Errorf("inCurfew(%02d:30, 22, 8) = %v, want %v", c.hour, got, c.want)
		}
	}
}

func TestInCurfewSameDayWindow(t *testing.T) {
	if !inCurfew(at(13), 9, 17) {
		t.Error("13:30 should be inside a 9-17 window")
	}
	if inCurfew(at(17), 9, 17) {
		t.Error("17:30 should be outside a 9-17 window (end is exclusive)")
	}
}
