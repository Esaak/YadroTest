package utils

import "time"

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}
