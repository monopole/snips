package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	DayFormat1 = "2006-01-02"
	DayFormat2 = "2006-Jan-02"
)

// DayRange is a specific calendar day (YYYY/MM/DD) paired with a day count.
// The day count is inclusive of the specific day.
// A day count less than one is illegal; there must be at least one day in the range.
type DayRange struct {
	Year     int
	Month    time.Month
	Day      int
	DayCount int
}

func MakeDayRange(dayStart string, dayCount int) (*DayRange, error) {
	if dayCount < 1 {
		return nil, fmt.Errorf("dayCount must be greater than zero")
	}

	var dayStartAsTime time.Time
	if dayStart == "" {
		// Default is today minus dayCount, revealing recent activity by default.
		dayStartAsTime = time.Now().Round(24*time.Hour).AddDate(0, 0, -dayCount)
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
		return fmt.Sprintf("%s - %s", d1.Format(f), d2.Format(f))
	}
	if d1.Month() != d2.Month() {
		const f = "January 2"
		return fmt.Sprintf("%s - %s %d", d1.Format(f), d2.Format(f), d2.Year())
	}
	if d1.Day() != d2.Day() {
		const f = "January 2"
		return fmt.Sprintf("%s-%d %d", d1.Format(f), d2.Day(), d2.Year())
	}
	return fmt.Sprintf("%s", d1.Format("January 2, 2006"))
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
