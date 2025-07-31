package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseMonthYear parses date string in format "MM-YYYY" to time.Time
// Returns first day of the month at 00:00:00
func ParseMonthYear(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, ErrInvalidDateFormat
	}

	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return time.Time{}, ErrInvalidDateFormat
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, ErrInvalidDateFormat
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil || year < 1900 || year > 3000 {
		return time.Time{}, ErrInvalidDateFormat
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

// FormatMonthYear formats time.Time to "MM-YYYY" string
func FormatMonthYear(t time.Time) string {
	return fmt.Sprintf("%02d-%04d", t.Month(), t.Year())
}

// CalculateMonthsInPeriod calculates the number of full months between two dates
func CalculateMonthsInPeriod(startDate, endDate time.Time) int {
	if endDate.Before(startDate) {
		return 0
	}

	years := endDate.Year() - startDate.Year()
	months := int(endDate.Month()) - int(startDate.Month())

	totalMonths := years*12 + months

	if endDate.Day() >= startDate.Day() {
		totalMonths++
	} else if totalMonths > 0 {
		totalMonths++
	}

	if totalMonths < 0 {
		return 0
	}

	return totalMonths
}

// CalculateSubscriptionMonthsInPeriod calculates how many months a subscription
// is active within the given period
func CalculateSubscriptionMonthsInPeriod(subStart, subEnd *time.Time, periodStart, periodEnd time.Time) int {
	actualStart := periodStart
	if subStart.After(periodStart) {
		actualStart = *subStart
	}

	actualEnd := periodEnd
	if subEnd != nil && subEnd.Before(periodEnd) {
		actualEnd = *subEnd
	}

	if actualStart.After(actualEnd) {
		return 0
	}

	return CalculateMonthsInPeriod(actualStart, actualEnd)
}

// GetLastDayOfMonth returns the last day of the month for given time
func GetLastDayOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 0, 23, 59, 59, 999999999, t.Location())
}
