package utils

import (
	"fmt"
	"time"
)

// JakartaLocation represents the Asia/Jakarta timezone
var JakartaLocation *time.Location

func init() {
	var err error
	JakartaLocation, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC+7 if timezone data is not available
		JakartaLocation = time.FixedZone("WIB", 7*60*60)
	}
}

// FormatDate formats a date according to the given format string
func FormatDate(t time.Time, format string) string {
	jakartaTime := t.In(JakartaLocation)

	switch format {
	case "yyyy-MM-dd":
		return jakartaTime.Format("2006-01-02")
	case "dd MMMM yyyy":
		return jakartaTime.Format("02 January 2006")
	case "dd/MM/yyyy":
		return jakartaTime.Format("02/01/2006")
	default:
		return jakartaTime.Format(format)
	}
}

// FormatTime formats a time according to the given format string
func FormatTime(t time.Time, format string) string {
	jakartaTime := t.In(JakartaLocation)

	switch format {
	case "HH:mm":
		return jakartaTime.Format("15:04")
	case "HH:mm:ss":
		return jakartaTime.Format("15:04:05")
	default:
		return jakartaTime.Format(format)
	}
}

// IsToday checks if the given time is today in Jakarta timezone
func IsToday(t time.Time) bool {
	now := time.Now().In(JakartaLocation)
	target := t.In(JakartaLocation)

	return now.Year() == target.Year() &&
		now.Month() == target.Month() &&
		now.Day() == target.Day()
}

// IsYesterday checks if the given time is yesterday in Jakarta timezone
func IsYesterday(t time.Time) bool {
	now := time.Now().In(JakartaLocation)
	yesterday := now.AddDate(0, 0, -1)
	target := t.In(JakartaLocation)

	return yesterday.Year() == target.Year() &&
		yesterday.Month() == target.Month() &&
		yesterday.Day() == target.Day()
}

// GetTodayDate returns today's date in YYYY-MM-DD format
func GetTodayDate() string {
	return FormatDate(time.Now(), "yyyy-MM-dd")
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// AddDays adds the specified number of days to the given time
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// CalculateWorkDuration calculates the duration between check-in and check-out times
func CalculateWorkDuration(checkIn, checkOut time.Time) string {
	duration := checkOut.Sub(checkIn)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d jam %d menit", hours, minutes)
	}
	return fmt.Sprintf("%d menit", minutes)
}

// NowInJakarta returns the current time in Jakarta timezone
func NowInJakarta() time.Time {
	return time.Now().In(JakartaLocation)
}
