package attendance

import (
	"attendance-bot/internal/database"
	"attendance-bot/internal/utils"
	"attendance-bot/pkg/models"
	"fmt"
	"strings"
	"time"
)

// Service handles attendance business logic
type Service struct {
	repo *database.Repository
	totp *TOTPService
}

// AttendanceResult represents the result of an attendance operation
type AttendanceResult struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Record  *models.AttendanceRecord `json:"record,omitempty"`
}

// NewService creates a new attendance service
func NewService(repo *database.Repository, totpSecret string) *Service {
	return &Service{
		repo: repo,
		totp: NewTOTPService(totpSecret),
	}
}

// MarkAttendance processes an attendance request
func (s *Service) MarkAttendance(userID int64, username, firstName string, lastName *string, otp string) (*AttendanceResult, error) {
	// Validate OTP
	if !utils.ValidateOTP(otp) {
		return &AttendanceResult{
			Success: false,
			Message: "❌ Format OTP tidak valid. Harap masukkan 6 digit angka.",
		}, nil
	}

	// Verify TOTP
	if !s.totp.Verify(otp) {
		return &AttendanceResult{
			Success: false,
			Message: "❌ Kode OTP tidak valid atau sudah kedaluwarsa. Silakan coba dengan kode yang baru.",
		}, nil
	}

	// Get current date and time
	now := utils.NowInJakarta()
	dateKey := utils.FormatDate(now, "yyyy-MM-dd")

	// Check current attendance status
	status, err := s.repo.GetUserAttendanceStatus(userID, dateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance status: %w", err)
	}

	// Determine attendance type and validate
	var attendanceType string
	var message string

	if !status.HasCheckedIn {
		// First attendance of the day - check in
		attendanceType = "check_in"
		timeStr := utils.FormatTime(now, "HH:mm")
		message = fmt.Sprintf("✅ **Absen Masuk** tercatat!\n⏰ Waktu: %s", timeStr)
	} else if !status.HasCheckedOut {
		// Second attendance of the day - check out
		attendanceType = "check_out"
		checkInTime := status.CheckInRecord.Timestamp
		timeStr := utils.FormatTime(now, "HH:mm")
		workDuration := utils.CalculateWorkDuration(checkInTime, now)
		message = fmt.Sprintf("🏠 **Absen Pulang** tercatat!\n⏰ Waktu: %s\n⌛ Durasi kerja: %s", timeStr, workDuration)
	} else {
		// Both check-in and check-out already done
		return &AttendanceResult{
			Success: false,
			Message: "❌ Anda sudah absen lengkap hari ini (masuk dan pulang)!",
		}, nil
	}

	// Create attendance record
	record := &models.AttendanceRecord{
		UserID:    userID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Timestamp: now,
		Type:      attendanceType,
		Date:      dateKey,
	}

	// Insert into database
	savedRecord, err := s.repo.InsertAttendance(record)
	if err != nil {
		return nil, fmt.Errorf("failed to save attendance: %w", err)
	}

	return &AttendanceResult{
		Success: true,
		Message: message,
		Record:  savedRecord,
	}, nil
}

// GetUserAttendanceStatus returns a user's attendance status for today
func (s *Service) GetUserAttendanceStatus(userID int64, date string) (*models.AttendanceStatus, error) {
	return s.repo.GetUserAttendanceStatus(userID, date)
}

// GetUserAttendanceHistory returns a user's attendance history
func (s *Service) GetUserAttendanceHistory(userID int64, days int) ([]models.AttendanceRecord, error) {
	return s.repo.GetUserAttendanceHistory(userID, days)
}

// GenerateAttendanceReport creates a formatted daily attendance report
func (s *Service) GenerateAttendanceReport() (string, error) {
	today := utils.GetTodayDate()
	records, err := s.repo.GetDailyReport(today)
	if err != nil {
		return "", fmt.Errorf("failed to get daily report: %w", err)
	}

	if len(records) == 0 {
		return "📭 Belum ada yang absen hari ini.", nil
	}

	// Group records by user to show check-in and check-out together
	userRecords := make(map[int64]map[string]*models.AttendanceRecord)
	for _, record := range records {
		if userRecords[record.UserID] == nil {
			userRecords[record.UserID] = make(map[string]*models.AttendanceRecord)
		}
		userRecords[record.UserID][record.Type] = &record
	}

	// Build report message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("📊 **Laporan Absensi Hari Ini**\n📅 %s\n\n",
		utils.FormatDate(time.Now(), "dd MMMM yyyy")))

	checkInCount := 0
	checkOutCount := 0
	userIndex := 1

	for _, userRecs := range userRecords {
		checkInRec := userRecs["check_in"]
		checkOutRec := userRecs["check_out"]

		if checkInRec != nil {
			name := s.formatUserName(checkInRec)
			checkInTime := utils.FormatTime(checkInRec.Timestamp, "HH:mm")

			message.WriteString(fmt.Sprintf("%d. **%s**\n", userIndex, name))
			message.WriteString(fmt.Sprintf("   ⏰ Masuk: %s", checkInTime))

			// Add status indicator for late arrival (after 9:00 AM)
			if checkInRec.Timestamp.Hour() >= 9 {
				message.WriteString(" ⚠️")
			} else {
				message.WriteString(" ✅")
			}
			message.WriteString("\n")

			checkInCount++
		}

		if checkOutRec != nil {
			if checkInRec == nil {
				// Handle edge case where there's check-out but no check-in
				name := s.formatUserName(checkOutRec)
				message.WriteString(fmt.Sprintf("%d. **%s**\n", userIndex, name))
				message.WriteString("   ⏰ Masuk: -\n")
			}

			checkOutTime := utils.FormatTime(checkOutRec.Timestamp, "HH:mm")
			message.WriteString(fmt.Sprintf("   🏠 Pulang: %s\n", checkOutTime))

			// Calculate work duration if both check-in and check-out exist
			if checkInRec != nil {
				duration := utils.CalculateWorkDuration(checkInRec.Timestamp, checkOutRec.Timestamp)
				message.WriteString(fmt.Sprintf("   ⌛ Durasi: %s\n", duration))
			}

			checkOutCount++
		} else if checkInRec != nil {
			message.WriteString("   🏠 Pulang: -\n")
		}

		message.WriteString("\n")
		userIndex++
	}

	// Add summary
	message.WriteString("**Ringkasan:**\n")
	message.WriteString(fmt.Sprintf("👥 Total Karyawan: %d\n", len(userRecords)))
	message.WriteString(fmt.Sprintf("📝 Check-in: %d\n", checkInCount))
	message.WriteString(fmt.Sprintf("🏠 Check-out: %d", checkOutCount))

	return message.String(), nil
}

// SetUserAlias sets a custom display name for a user
func (s *Service) SetUserAlias(userID int64, firstName string, lastName *string) error {
	return s.repo.SetUserAlias(userID, firstName, lastName)
}

// formatUserName returns the display name for a user, preferring alias if available
func (s *Service) formatUserName(record *models.AttendanceRecord) string {
	// Try to get alias first
	alias, err := s.repo.GetUserAlias(record.UserID)
	if err == nil && alias != nil {
		if alias.LastName != nil && *alias.LastName != "" {
			return fmt.Sprintf("%s %s", alias.FirstName, *alias.LastName)
		}
		return alias.FirstName
	}

	// Fall back to original name
	if record.LastName != nil && *record.LastName != "" {
		return fmt.Sprintf("%s %s", record.FirstName, *record.LastName)
	}
	return record.FirstName
}

// GetAttendanceReportRange generates a report for a date range
func (s *Service) GetAttendanceReportRange(startDate, endDate string) ([]models.AttendanceRecord, error) {
	return s.repo.GetAttendanceReportRange(startDate, endDate)
}
