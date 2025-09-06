package reports

import (
	"attendance-bot/internal/utils"
	"attendance-bot/pkg/models"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CSVGenerator handles CSV report generation
type CSVGenerator struct {
	outputDir string
}

// NewCSVGenerator creates a new CSV generator
func NewCSVGenerator(outputDir string) *CSVGenerator {
	return &CSVGenerator{
		outputDir: outputDir,
	}
}

// GenerateAttendanceReport creates a CSV file with attendance data
func (g *CSVGenerator) GenerateAttendanceReport(records []models.AttendanceRecord, startDate, endDate string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("attendance_report_%s_to_%s.csv", startDate, endDate)
	filepath := filepath.Join(g.outputDir, filename)

	// Create CSV file
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"User ID",
		"Username",
		"First Name",
		"Last Name",
		"Date",
		"Type",
		"Time",
		"Timestamp",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write records
	for _, record := range records {
		lastName := ""
		if record.LastName != nil {
			lastName = *record.LastName
		}

		timeStr := utils.FormatTime(record.Timestamp, "HH:mm:ss")

		row := []string{
			fmt.Sprintf("%d", record.ID),
			fmt.Sprintf("%d", record.UserID),
			record.Username,
			record.FirstName,
			lastName,
			record.Date,
			record.Type,
			timeStr,
			record.Timestamp.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return filepath, nil
}

// GenerateDailyReport creates a CSV for a specific date
func (g *CSVGenerator) GenerateDailyReport(records []models.AttendanceRecord, date string) (string, error) {
	return g.GenerateAttendanceReport(records, date, date)
}

// GenerateUserReport creates a CSV for a specific user's attendance
func (g *CSVGenerator) GenerateUserReport(records []models.AttendanceRecord, userID int64, days int) (string, error) {
	if len(records) == 0 {
		return "", fmt.Errorf("no records found for user %d", userID)
	}

	// Use the date range from the records
	startDate := records[len(records)-1].Date // oldest
	endDate := records[0].Date                // newest

	filename := fmt.Sprintf("user_%d_attendance_%s_to_%s.csv", userID, startDate, endDate)
	filepath := filepath.Join(g.outputDir, filename)

	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create CSV file
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Date",
		"Check-in Time",
		"Check-out Time",
		"Work Duration",
		"Status",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Group records by date
	dailyRecords := make(map[string]map[string]*models.AttendanceRecord)
	for _, record := range records {
		if dailyRecords[record.Date] == nil {
			dailyRecords[record.Date] = make(map[string]*models.AttendanceRecord)
		}
		dailyRecords[record.Date][record.Type] = &record
	}

	// Write records grouped by date
	for date, dayRecords := range dailyRecords {
		checkIn := dayRecords["check_in"]
		checkOut := dayRecords["check_out"]

		checkInTime := "-"
		checkOutTime := "-"
		duration := "-"
		status := "Absent"

		if checkIn != nil {
			checkInTime = utils.FormatTime(checkIn.Timestamp, "HH:mm:ss")
			status = "Present"
			if checkIn.Timestamp.Hour() >= 9 {
				status = "Late"
			}
		}

		if checkOut != nil {
			checkOutTime = utils.FormatTime(checkOut.Timestamp, "HH:mm:ss")
			if checkIn != nil {
				duration = utils.CalculateWorkDuration(checkIn.Timestamp, checkOut.Timestamp)
			}
		}

		row := []string{
			date,
			checkInTime,
			checkOutTime,
			duration,
			status,
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return filepath, nil
}
