package models

import "time"

// AttendanceRecord represents a single attendance entry
type AttendanceRecord struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  *string   `json:"last_name,omitempty" db:"last_name"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Type      string    `json:"type" db:"type"` // "check_in" or "check_out"
	Date      string    `json:"date" db:"date"` // YYYY-MM-DD format
}

// UserAlias represents a user's custom display name
type UserAlias struct {
	UserID    int64   `json:"user_id" db:"user_id"`
	FirstName string  `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name,omitempty" db:"last_name"`
}

// AttendanceStatus represents a user's attendance status for a given day
type AttendanceStatus struct {
	HasCheckedIn   bool              `json:"has_checked_in"`
	HasCheckedOut  bool              `json:"has_checked_out"`
	CheckInRecord  *AttendanceRecord `json:"check_in_record,omitempty"`
	CheckOutRecord *AttendanceRecord `json:"check_out_record,omitempty"`
}
