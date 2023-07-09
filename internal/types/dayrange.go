package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	DayFormatGitHub     = "2006-01-02"
	DayFormatHuman      = "2006-Jan-02"
	DayFormatJira       = "2006/01/02"
	DateFormatJiraIssue = "2006-01-02T15:04:05.000-0700"
	defaultDayCount     = 14 // two weeks
)

// DayRange is a specific calendar start day (year, month, day number) paired with a DayCount.
// The end date of a DayRange is the start date plus (DayCount - 1).
// DayRange intentionally doesn't use an instance of time.Time to eliminate confusion on what
// to do with hours, minutes, seconds, etc.
// A DayCount less than one is illegal; there must be at least one day in the range.
// Handy for GitHub date range queries.
// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
type DayRange struct {
	Year     int
	Month    time.Month
	Day      int
	DayCount int
}

// MakeDayRange makes an instance of DayRange from the given arguments.
func MakeDayRange(dayStart string, dayEnd string, dayCount int) (*DayRange, error) {
	if dayEnd != "" {
		return makeDayRangeWithExplicitEndDate(dayStart, dayEnd, dayCount)
	}
	if dayStart == "" {
		// Neither the start nor the end day was specified.
		// Use "today" as end date, and count backwards with dayCount.
		if dayCount < 1 {
			dayCount = defaultDayCount
		}
		return makeDayRangeFromEnd(today(), dayCount), nil
	}
	// Only dayStart was specified.
	start, err := parseDate(dayStart)
	if err != nil {
		return nil, err
	}
	if dayCount < 1 {
		end := today()
		if start.After(end) {
			return nil, fmt.Errorf("start date %s is after today (%s)", dayStart, end.Format(DayFormatHuman))
		}
		dayCount = dayCountInclusive(start, end)
	}
	return makeDayRangeFromStart(start, dayCount), nil
}

func makeDayRangeWithExplicitEndDate(dayStart string, dayEnd string, dayCount int) (*DayRange, error) {
	if dayEnd == "" {
		return nil, fmt.Errorf("this code path wants a non-empty dayEnd")
	}
	var start, end time.Time
	var err error
	end, err = parseDate(dayEnd)
	if err != nil {
		return nil, err
	}
	if dayStart == "" {
		if dayCount < 1 {
			dayCount = defaultDayCount
		}
		return makeDayRangeFromEnd(end, dayCount), nil
	}
	if dayCount > 0 {
		return nil, fmt.Errorf("specify only 2 of dayStart, dayEnd and dayCount")
	}
	start, err = parseDate(dayStart)
	if err != nil {
		return nil, err
	}
	if !start.Before(end) {
		return nil, fmt.Errorf("dayStart must preceed dayEnd")
	}
	hours := int(end.Sub(start).Hours())
	if hours%24 != 0 {
		return nil, fmt.Errorf("hours=%d not divisible by 24, wtf", hours)
	}
	dayCount = (hours / 24) + 1
	return makeDayRangeFromStart(start, dayCount), nil
}

// today returns a timestamp matching today.
func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func makeDayRangeFromEnd(end time.Time, dayCount int) *DayRange {
	// dayCount is an inclusive count, so if the start and end date are the same,
	// dayCount is one, not zero.
	return makeDayRangeFromStart(end.AddDate(0, 0, -(dayCount-1)), dayCount)
}

func makeDayRangeFromStart(start time.Time, c int) *DayRange {
	return &DayRange{
		Year:     start.Year(),
		Month:    start.Month(),
		Day:      start.Day(),
		DayCount: c,
	}
}

func dayCountInclusive(t0, t1 time.Time) int {
	diff := t1.Sub(t0).Truncate(24 * time.Hour)
	days := int(diff.Hours()) / 24
	return days + 1
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
	return []string{DayFormatGitHub, DayFormatHuman, DayFormatJira}
}

func DateOptions() string {
	opts := AllDateFormats()
	return strings.Join(opts[0:len(opts)-1], ", ") + " or " + opts[len(opts)-1]
}
