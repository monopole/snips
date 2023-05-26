package internal

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const dayFormat = "2006-01-02"

func ParseDate(arg string) time.Time {
	t, err := time.Parse(dayFormat, arg)
	if err != nil {
		fmt.Printf("Trouble with date specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return t
}

func ParseDayCount(arg string) int {
	i, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("Trouble with day count specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return i
}

func DateRange(d1, d2 time.Time) string {
	if d1.Year() != d2.Year() {
		return fmt.Sprintf(
			"%s - %s",
			d1.Format("January 2, 2006"),
			d2.Format("January 2, 2006"))
	}
	if d1.Month() != d2.Month() {
		return fmt.Sprintf(
			"%s - %s %d",
			d1.Format("January 2"),
			d2.Format("January 2"), d2.Year())
	}
	return fmt.Sprintf(
		"%s %d-%d, %d",
		d1.Format("January"),
		d1.Day(), d2.Day(), d2.Year())
}
