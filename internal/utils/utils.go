package utils

import "time"

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}
