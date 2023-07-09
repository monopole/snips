package types

import (
	"testing"
	"time"
)

func Test_dayCountInclusive(t *testing.T) {
	type testCase struct {
		dayStart        string
		dayEnd          string
		desiredDayCount int
	}
	tests := map[string]testCase{
		"t1": {
			dayStart:        "2020-Mar-18",
			dayEnd:          "2020-Mar-18",
			desiredDayCount: 1,
		},
		"t2": {
			dayStart:        "2020-Mar-18",
			dayEnd:          "2020-Mar-19",
			desiredDayCount: 2,
		},
		"t3": {
			dayStart:        "2023-May-31",
			dayEnd:          "2023-Jun-03",
			desiredDayCount: 4,
		},
		"t4": {
			dayStart:        "2023-Jan-01",
			dayEnd:          "2023-Jun-23",
			desiredDayCount: 174,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t0, err := parseDate(tc.dayStart)
			if err != nil {
				t.Error(err.Error())
			}
			t1, err := parseDate(tc.dayEnd)
			if err != nil {
				t.Error(err.Error())
			}
			dc := dayCountInclusive(t0, t1)
			if dc != tc.desiredDayCount {
				t.Fatalf("expected dc=%d, got %d", tc.desiredDayCount, dc)
			}
		})
	}

}

func TestDayRange_MakeDayRange(t *testing.T) {
	type testCase struct {
		dayStart string
		dayEnd   string
		dayCount int
		want     string
	}
	tests := map[string]testCase{
		"t1": {
			dayStart: "2020-Mar-18",
			dayCount: 1,
			want:     "March 18, 2020 (one day)",
		},
		"t2": {
			dayStart: "2020-Mar-01",
			dayCount: 1,
			want:     "March 1, 2020 (one day)",
		},
		"t3": {
			dayStart: "2020-03-01",
			dayCount: 1,
			want:     "March 1, 2020 (one day)",
		},
		"t4": {
			dayStart: "2020-Mar-30",
			dayCount: 5, // March has 31 days
			want:     "March 30 - April 3 2020 (5 days)",
		},
		"t5": {
			dayStart: "2020-Mar-30",
			dayEnd:   "2020-Apr-03",
			want:     "March 30 - April 3 2020 (5 days)",
		},
		"t6": {
			dayStart: "2020-Dec-30",
			dayCount: 5, // December has 31 days
			want:     "December 30, 2020 - January 3, 2021 (5 days)",
		},
		"t7": {
			dayStart: "2020-Dec-30",
			dayEnd:   "2021-Jan-03",
			want:     "December 30, 2020 - January 3, 2021 (5 days)",
		},
		"t8": {
			dayStart: "2023-Jan-01",
			dayEnd:   "2023-Jun-23",
			want:     "January 1 - June 23 2023 (174 days)",
		},
		"t9": {
			dayStart: "",
			dayEnd:   "2023-Jun-23",
			want:     "June 10-23 2023 (14 days)",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dr, err := MakeDayRange(tc.dayStart, tc.dayEnd, tc.dayCount)
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
			if got := dr.StartAsTime().Format(DayFormatGitHub); got != tc.want {
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
			if got := dr.EndAsTime().Format(DayFormatGitHub); got != tc.want {
				t.Errorf("EndAsTime() = %v, want %v", got, tc.want)
			}
		})
	}
}
