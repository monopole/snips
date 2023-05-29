package types

import (
	"testing"
	"time"
)

func TestDayRange_MakeDayRange(t *testing.T) {
	type testCase struct {
		arg      string
		dayCount int
		want     string
	}
	tests := map[string]testCase{
		"t1": {
			arg:      "2020-Mar-18",
			dayCount: 1,
			want:     "March 18, 2020",
		},
		"t2": {
			arg:      "2020-Mar-30",
			dayCount: 5, // March has 31 days
			want:     "March 30 - April 3 2020",
		},
		"t3": {
			arg:      "2020-Dec-30",
			dayCount: 5, // December has 31 days
			want:     "December 30, 2020 - January 3, 2021",
		},
		"t4": {
			arg:      "2020-Mar-01",
			dayCount: 1,
			want:     "March 1, 2020",
		},
		"t5": {
			arg:      "2020-03-01",
			dayCount: 1,
			want:     "March 1, 2020",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dr, err := MakeDayRange(tc.arg, tc.dayCount)
			if err != nil {
				t.Fatalf(err.Error())
			}
			got := dr.PrettyRange()
			if got != tc.want {
				t.Errorf("MakeDayRange() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDayRange_StartAsTime(t *testing.T) {
	type testCase struct {
		y    int
		m    time.Month
		d    int
		want string
	}
	tests := map[string]testCase{
		"t1": {
			y:    2020,
			m:    3,
			d:    18,
			want: "2020-03-18",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dr := &DayRange{
				Year:  tc.y,
				Month: tc.m,
				Day:   tc.d,
			}
			if got := dr.StartAsTime().Format(DayFormat1); got != tc.want {
				t.Errorf("StartAsTime() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDayRange_EndAsTime(t *testing.T) {
	type testCase struct {
		y        int
		m        time.Month
		d        int
		dayCount int
		want     string
	}
	tests := map[string]testCase{
		"t1": {
			y:        2020,
			m:        3,
			d:        18,
			dayCount: 1,
			want:     "2020-03-18",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dr := &DayRange{
				Year:     tc.y,
				Month:    tc.m,
				Day:      tc.d,
				DayCount: tc.dayCount,
			}
			if got := dr.EndAsTime().Format(DayFormat1); got != tc.want {
				t.Errorf("EndAsTime() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDayRange_PrettyRange(t *testing.T) {
	type testCase struct {
		y        int
		m        time.Month
		d        int
		dayCount int
		want     string
	}
	tests := map[string]testCase{
		"t1": {
			y:        2020,
			m:        3,
			d:        18,
			dayCount: 1,
			want:     "March 18, 2020",
		},
		"t2": {
			y:        2020,
			m:        3,
			d:        30, // March has 31 days
			dayCount: 5,
			want:     "March 30 - April 3 2020",
		},
		"t3": {
			y:        2020,
			m:        12,
			d:        30, // December has 31 days
			dayCount: 5,
			want:     "December 30, 2020 - January 3, 2021",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dr := &DayRange{
				Year:     tc.y,
				Month:    tc.m,
				Day:      tc.d,
				DayCount: tc.dayCount,
			}
			if got := dr.PrettyRange(); got != tc.want {
				t.Errorf("func TestDayRange_PrettyRange(t *testing.T) {\n() = %v, want %v", got, tc.want)
			}
		})
	}
}
