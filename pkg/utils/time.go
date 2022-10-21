package utils

import (
	"math"
	"time"
)

func DurationProcess() func() time.Duration {
	startTime := time.Now()

	return func() time.Duration {
		endTime := time.Now()
		return endTime.Sub(startTime)
	}
}

func DiffDays(startDate time.Time, endDate time.Time) int {
	diffTime := endDate.Sub(startDate)
	diffDays := math.Round(diffTime.Hours() / 24)
	return int(diffDays)
}

func ToUTC(date time.Time, loc *time.Location) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), 0, loc).UTC()
}

func InDateRange(start, end, compare time.Time) bool {
	if start.After(compare) || end.Before(compare) {
		return false
	}

	return true
}

func IsSameDate(one, compare time.Time) bool {
	y1, m1, d1 := one.Date()
	y2, m2, d2 := compare.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		return true
	}

	return false
}

func IsTheLastDayOfMonth(date time.Time) bool {
	return date.Month() != date.AddDate(0, 0, 1).Month()
}