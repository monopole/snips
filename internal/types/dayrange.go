package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	DayFormat1      = "2006-01-02"
	DayFormat2      = "2006-Jan-02"
	defaultDayCount = 14 // two weeks
)

// DayRange is a specific calendar day (YYYY/MM/DD) paired with a day count.
// The day count is inclusive of the specific day.
// A day count less than one is illegal; there must be at least one day in the range.
// The end calendar day is computed from the y/m/d start date, and used in
// GitHub range queries:
// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
type DayRange struct {
	Year     int
	Month    time.Month
	Day      int
	DayCount int
}

func makeDayRangeWithExplicitEndDate(dayStart string, dayEnd string, dayCount int) (*DayRange, error) {
	if dayEnd == "" {
		return nil, fmt.Errorf("this code path wants a non-empty dayEnd")
	}
	var dayStartAsTime, dayEndAsTime time.Time
	var err error
	dayEndAsTime, err = parseDate(dayEnd)
	if err != nil {
		return nil, err
	}
	if dayStart == "" {
		if dayCount < 1 {
			dayCount = defaultDayCount
		}
		dayStartAsTime = dayEndAsTime.AddDate(0, 0, -dayCount)
	} else {
		if dayCount > 0 {
			return nil, fmt.Errorf("specify only 2 of dayStart, dayEnd and dayCount")
		}
		dayStartAsTime, err = parseDate(dayStart)
		if err != nil {
			return nil, err
		}
		if !dayStartAsTime.Before(dayEndAsTime) {
			return nil, fmt.Errorf("dayStart must preceed dayEnd")
		}
		hours := int(dayEndAsTime.Sub(dayStartAsTime).Hours())
		if hours%24 != 0 {
			return nil, fmt.Errorf("hours=%d not divisible by 24, wtf", hours)
		}
		dayCount = (hours / 24) + 1
	}
	return &DayRange{
		Year:     dayStartAsTime.Year(),
		Month:    dayStartAsTime.Month(),
		Day:      dayStartAsTime.Day(),
		DayCount: dayCount,
	}, nil
}

// MakeDayRange makes an instance of DayRange from the given arguments.
// Date strings must be in the format YYYY/MM/DD.
func MakeDayRange(dayStart string, dayEnd string, dayCount int) (*DayRange, error) {
	if dayEnd != "" {
		return makeDayRangeWithExplicitEndDate(dayStart, dayEnd, dayCount)
	}
	var dayStartAsTime time.Time
	if dayCount < 1 {
		dayCount = defaultDayCount
	}
	if dayStart == "" {
		// Default is today minus dayCount-1, revealing recent activity by default.
		// Subtracting 1 to assure that "today" is in the range.
		dayStartAsTime = time.Now().Round(24*time.Hour).AddDate(0, 0, -(dayCount - 1))
	} else {
		var err error
		dayStartAsTime, err = parseDate(dayStart)
		if err != nil {
			return nil, err
		}
	}
	return &DayRange{
		Year:     dayStartAsTime.Year(),
		Month:    dayStartAsTime.Month(),
		Day:      dayStartAsTime.Day(),
		DayCount: dayCount,
	}, nil
}

func (dr *DayRange) StartAsTime() time.Time {
	return time.Date(dr.Year, dr.Month, dr.Day, 0, 0, 0, 0, time.Local)
}

func (dr *DayRange) EndAsTime() time.Time {
	// Recall that dayCount must be >= 1.
	return dr.StartAsTime().AddDate(0, 0, dr.DayCount-1)
}

// PrettyRange returns a simplified date range as a string.
func (dr *DayRange) PrettyRange() string {
	d1 := dr.StartAsTime()
	d2 := dr.EndAsTime()
	if d1.Year() != d2.Year() {
		const f = "January 2, 2006"
		return fmt.Sprintf("%s - %s (%d days)", d1.Format(f), d2.Format(f), dr.DayCount)
	}
	if d1.Month() != d2.Month() {
		const f = "January 2"
		return fmt.Sprintf("%s - %s %d (%d days)", d1.Format(f), d2.Format(f), d2.Year(), dr.DayCount)
	}
	if d1.Day() != d2.Day() {
		const f = "January 2"
		return fmt.Sprintf("%s-%d %d (%d days)", d1.Format(f), d2.Day(), d2.Year(), dr.DayCount)
	}
	return fmt.Sprintf("%s (one day)", d1.Format("January 2, 2006"))
}

func parseDate(v string) (time.Time, error) {
	for _, f := range AllDateFormats() {
		if t, err := time.Parse(f, v); err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("bad date value %q, use formats %s", v, DateOptions())
}

func AllDateFormats() []string {
	return []string{DayFormat1, DayFormat2}
}

func DateOptions() string {
	opts := AllDateFormats()
	return strings.Join(opts[0:len(opts)-1], ", ") + " or " + opts[len(opts)-1]
}
