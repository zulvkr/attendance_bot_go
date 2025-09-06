package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateOTP checks if the provided string is a valid 6-digit OTP
func ValidateOTP(otp string) bool {
	// Remove any whitespace
	otp = strings.TrimSpace(otp)

	// Check if it's exactly 6 digits
	if len(otp) != 6 {
		return false
	}

	// Check if all characters are digits
	matched, err := regexp.MatchString(`^\d{6}$`, otp)
	if err != nil {
		return false
	}

	return matched
}

// IsValidTelegramUserID checks if the provided user ID is valid
func IsValidTelegramUserID(userID int64) bool {
	return userID > 0
}

// SanitizeUsername removes potentially harmful characters from username
func SanitizeUsername(username string) string {
	// Remove any non-alphanumeric characters except underscore and hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	return reg.ReplaceAllString(username, "")
}

// SanitizeName removes potentially harmful characters from names
func SanitizeName(name string) string {
	// Allow letters, spaces, apostrophes, and hyphens
	reg := regexp.MustCompile(`[^a-zA-Z\s'\-]`)
	cleaned := reg.ReplaceAllString(name, "")

	// Trim whitespace and limit length
	cleaned = strings.TrimSpace(cleaned)
	if len(cleaned) > 50 {
		cleaned = cleaned[:50]
	}

	return cleaned
}

// ParseInteger safely parses a string to integer
func ParseInteger(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

// IsValidDateFormat checks if the date is in YYYY-MM-DD format
func IsValidDateFormat(date string) bool {
	matched, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	if err != nil {
		return false
	}
	return matched
}
