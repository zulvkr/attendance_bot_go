package database

import (
	"attendance-bot/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// Repository handles all database operations
type Repository struct {
	db *SQLiteDB
}

// NewRepository creates a new repository instance
func NewRepository(db *SQLiteDB) *Repository {
	return &Repository{db: db}
}

// InsertAttendance adds a new attendance record
func (r *Repository) InsertAttendance(record *models.AttendanceRecord) (*models.AttendanceRecord, error) {
	query := `
		INSERT INTO attendance (user_id, username, first_name, last_name, timestamp, type, date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		record.UserID,
		record.Username,
		record.FirstName,
		record.LastName,
		record.Timestamp.Format(time.RFC3339),
		record.Type,
		record.Date,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert attendance: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	record.ID = id
	return record, nil
}

// GetUserAttendanceToday retrieves today's attendance records for a user
func (r *Repository) GetUserAttendanceToday(userID int64, date string) ([]models.AttendanceRecord, error) {
	query := `
		SELECT id, user_id, username, first_name, last_name, timestamp, type, date
		FROM attendance
		WHERE user_id = ? AND date = ?
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(query, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance: %w", err)
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		record, err := r.scanAttendanceRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

// GetUserAttendanceStatus returns the attendance status for a user on a specific date
func (r *Repository) GetUserAttendanceStatus(userID int64, date string) (*models.AttendanceStatus, error) {
	records, err := r.GetUserAttendanceToday(userID, date)
	if err != nil {
		return nil, err
	}

	status := &models.AttendanceStatus{
		HasCheckedIn:  false,
		HasCheckedOut: false,
	}

	for _, record := range records {
		if record.Type == "check_in" {
			status.HasCheckedIn = true
			status.CheckInRecord = &record
		} else if record.Type == "check_out" {
			status.HasCheckedOut = true
			status.CheckOutRecord = &record
		}
	}

	return status, nil
}

// GetUserAttendanceHistory retrieves attendance history for a user
func (r *Repository) GetUserAttendanceHistory(userID int64, days int) ([]models.AttendanceRecord, error) {
	query := `
		SELECT id, user_id, username, first_name, last_name, timestamp, type, date
		FROM attendance
		WHERE user_id = ? AND date >= date('now', '-' || ? || ' days')
		ORDER BY date DESC, timestamp ASC
	`

	rows, err := r.db.Query(query, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance history: %w", err)
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		record, err := r.scanAttendanceRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

// GetDailyReport retrieves all attendance records for a specific date
func (r *Repository) GetDailyReport(date string) ([]models.AttendanceRecord, error) {
	query := `
		SELECT a.id, a.user_id, a.username, a.first_name, a.last_name, a.timestamp, a.type, a.date
		FROM attendance a
		LEFT JOIN alias al ON a.user_id = al.user_id
		WHERE a.date = ?
		ORDER BY a.timestamp ASC
	`

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily report: %w", err)
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		record, err := r.scanAttendanceRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

// GetAttendanceReportRange retrieves attendance records within a date range
func (r *Repository) GetAttendanceReportRange(startDate, endDate string) ([]models.AttendanceRecord, error) {
	query := `
		SELECT a.id, a.user_id, a.username, a.first_name, a.last_name, a.timestamp, a.type, a.date
		FROM attendance a
		LEFT JOIN alias al ON a.user_id = al.user_id
		WHERE a.date BETWEEN ? AND ?
		ORDER BY a.date ASC, a.timestamp ASC
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance report range: %w", err)
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		record, err := r.scanAttendanceRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

// SetUserAlias sets or updates a user's alias
func (r *Repository) SetUserAlias(userID int64, firstName string, lastName *string) error {
	// Check if alias already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM alias WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existing alias: %w", err)
	}

	var query string
	var args []interface{}

	if exists {
		query = "UPDATE alias SET first_name = ?, last_name = ? WHERE user_id = ?"
		args = []interface{}{firstName, lastName, userID}
	} else {
		query = "INSERT INTO alias (user_id, first_name, last_name) VALUES (?, ?, ?)"
		args = []interface{}{userID, firstName, lastName}
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to set user alias: %w", err)
	}

	return nil
}

// GetUserAlias retrieves a user's alias
func (r *Repository) GetUserAlias(userID int64) (*models.UserAlias, error) {
	query := "SELECT user_id, first_name, last_name FROM alias WHERE user_id = ?"

	var alias models.UserAlias
	var lastName sql.NullString

	err := r.db.QueryRow(query, userID).Scan(&alias.UserID, &alias.FirstName, &lastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No alias found
		}
		return nil, fmt.Errorf("failed to get user alias: %w", err)
	}

	if lastName.Valid {
		alias.LastName = &lastName.String
	}

	return &alias, nil
}

// scanAttendanceRecord scans a database row into an AttendanceRecord
func (r *Repository) scanAttendanceRecord(rows *sql.Rows) (*models.AttendanceRecord, error) {
	var record models.AttendanceRecord
	var lastName sql.NullString
	var timestampStr string

	err := rows.Scan(
		&record.ID,
		&record.UserID,
		&record.Username,
		&record.FirstName,
		&lastName,
		&timestampStr,
		&record.Type,
		&record.Date,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan attendance record: %w", err)
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	record.Timestamp = timestamp

	// Handle nullable last name
	if lastName.Valid {
		record.LastName = &lastName.String
	}

	return &record, nil
}

// CheckUserAttendanceExists checks if a user has any attendance record for a specific date and type
func (r *Repository) CheckUserAttendanceExists(userID int64, date, attendanceType string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM attendance WHERE user_id = ? AND date = ? AND type = ?)"

	var exists bool
	err := r.db.QueryRow(query, userID, date, attendanceType).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check attendance existence: %w", err)
	}

	return exists, nil
}
