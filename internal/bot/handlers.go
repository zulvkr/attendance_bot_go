package bot

import (
	"attendance-bot/internal/attendance"
	"attendance-bot/internal/config"
	"attendance-bot/internal/reports"
	"attendance-bot/internal/utils"
	"attendance-bot/pkg/models"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"
)

// SessionData represents user session state
type SessionData struct {
	AwaitingDateRange bool
}

// Bot represents the main bot instance
type Bot struct {
	api               *TelegramAPI
	attendanceService *attendance.Service
	csvGenerator      *reports.CSVGenerator
	config            *config.Config
	logger            *slog.Logger
	lastUpdateID      int64
	sessions          map[int64]*SessionData // Simple in-memory session storage
}

// NewBot creates a new bot instance
func NewBot(token string, attendanceService *attendance.Service, csvGenerator *reports.CSVGenerator, cfg *config.Config, logger *slog.Logger) *Bot {
	return &Bot{
		api:               NewTelegramAPI(token),
		attendanceService: attendanceService,
		csvGenerator:      csvGenerator,
		config:            cfg,
		logger:            logger,
		sessions:          make(map[int64]*SessionData),
	}
}

// Start begins the bot polling loop
func (b *Bot) Start() error {
	b.logger.Info("Starting bot...")

	// Get bot info
	botInfo, err := b.api.GetMe()
	if err != nil {
		return fmt.Errorf("failed to get bot info: %w", err)
	}

	b.logger.Info("Bot started successfully", "bot_username", botInfo.Username, "bot_id", botInfo.ID)

	// Start polling loop
	for {
		updates, err := b.api.GetUpdates(b.lastUpdateID+1, 60)
		if err != nil {
			b.logger.Error("Failed to get updates", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			b.lastUpdateID = update.UpdateID
			if err := b.handleUpdate(&update); err != nil {
				b.logger.Error("Failed to handle update", "error", err, "update_id", update.UpdateID)
			}
		}
	}
}

// handleUpdate processes a single update
func (b *Bot) handleUpdate(update *Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message
	b.logger.Debug("Received message",
		"user_id", msg.From.ID,
		"username", msg.From.Username,
		"text", msg.Text)

	// Handle commands
	if strings.HasPrefix(msg.Text, "/") {
		return b.handleCommand(msg)
	}

	// Handle OTP (6-digit numbers)
	if utils.ValidateOTP(msg.Text) {
		return b.handleOTP(msg)
	}

	// Handle other text messages
	return b.handleTextMessage(msg)
}

// handleCommand processes bot commands
func (b *Bot) handleCommand(msg *Message) error {
	parts := strings.Fields(msg.Text)
	if len(parts) == 0 {
		return nil
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "/start":
		return b.handleStart(msg)
	case "/help":
		return b.handleHelp(msg)
	case "/report":
		return b.handleReport(msg)
	case "/history":
		return b.handleHistory(msg)
	case "/status":
		return b.handleStatus(msg)
	case "/alias":
		return b.handleAlias(msg, args)
	case "/fullreport":
		return b.handleFullReport(msg, args)
	default:
		return b.sendMessage(msg.Chat.ID, "❓ Perintah tidak dikenal. Ketik /help untuk melihat daftar perintah.")
	}
}

// handleStart handles the /start command
func (b *Bot) handleStart(msg *Message) error {
	welcomeMessage := `🎯 *Selamat datang di Attendance Bot!*

Untuk absen, kirimkan kode OTP 6 digit Anda.

*Perintah yang Tersedia:*
📝 Kirim OTP - Absen (masuk/pulang)
📊 /report - Lihat laporan absensi hari ini
📈 /history - Lihat riwayat absensi Anda
🏷️ /alias - Absen dengan nama lain
🔄 /status - Cek status absensi hari ini
📋 /fullreport - Download laporan lengkap (CSV)
❓ /help - Tampilkan pesan bantuan ini

*Sistem Absensi:*
• Absen pertama = Masuk (check-in)
• Absen kedua = Pulang (check-out)`

	return b.sendMarkdownMessage(msg.Chat.ID, welcomeMessage)
}

// handleHelp handles the /help command
func (b *Bot) handleHelp(msg *Message) error {
	helpMessage := `❓ *Bantuan Attendance Bot*

*Cara menggunakan:*
1. Dapatkan OTP dari aplikasi autentikator Anda
2. Kirimkan kode 6 digit ke bot ini
3. Sistem akan otomatis menentukan check-in atau check-out

*Sistem Absensi:*
• Absen pertama dalam hari = *Check-in* (Masuk)
• Absen kedua dalam hari = *Check-out* (Pulang)

*Perintah:*
📊 /report - Lihat laporan absensi hari ini
📈 /history - Lihat riwayat absensi Anda (30 hari terakhir)
🔄 /status - Cek status absensi hari ini (masuk/pulang)
🏷️ /alias - Gunakan nama panggilan/alias untuk absensi
   Format: /alias [Nama Depan] [Nama Belakang]
   Contoh: /alias John Doe
📋 /fullreport - Download laporan lengkap dalam format CSV
   Format: Masukkan rentang tanggal (YYYY-MM-DD YYYY-MM-DD)`

	return b.sendMarkdownMessage(msg.Chat.ID, helpMessage)
}

// handleReport handles the /report command
func (b *Bot) handleReport(msg *Message) error {
	report, err := b.attendanceService.GenerateAttendanceReport()
	if err != nil {
		b.logger.Error("Failed to generate report", "error", err)
		return b.sendMessage(msg.Chat.ID, "❌ Terjadi kesalahan saat membuat laporan. Silakan coba lagi.")
	}

	return b.sendMarkdownMessage(msg.Chat.ID, report)
}

// handleHistory handles the /history command
func (b *Bot) handleHistory(msg *Message) error {
	records, err := b.attendanceService.GetUserAttendanceHistory(msg.From.ID, 30)
	if err != nil {
		b.logger.Error("Failed to get attendance history", "error", err, "user_id", msg.From.ID)
		return b.sendMessage(msg.Chat.ID, "❌ Terjadi kesalahan saat mengambil riwayat. Silakan coba lagi.")
	}

	if len(records) == 0 {
		return b.sendMessage(msg.Chat.ID, "📭 Tidak ada riwayat absensi dalam 30 hari terakhir.")
	}

	message := b.formatHistoryMessage(records)
	return b.sendMarkdownMessage(msg.Chat.ID, message)
}

// handleStatus handles the /status command
func (b *Bot) handleStatus(msg *Message) error {
	today := utils.GetTodayDate()
	status, err := b.attendanceService.GetUserAttendanceStatus(msg.From.ID, today)
	if err != nil {
		b.logger.Error("Failed to get attendance status", "error", err, "user_id", msg.From.ID)
		return b.sendMessage(msg.Chat.ID, "❌ Terjadi kesalahan saat mengecek status. Silakan coba lagi.")
	}

	var message string
	if !status.HasCheckedIn && !status.HasCheckedOut {
		message = "❌ *Status Absensi*\n\nAnda belum absen hari ini.\nKirim OTP Anda untuk *check-in*."
	} else if status.HasCheckedIn && !status.HasCheckedOut {
		checkInTime := utils.FormatTime(status.CheckInRecord.Timestamp, "HH:mm")
		message = fmt.Sprintf("🟡 *Status Absensi*\n\n✅ Check-in: %s\n❌ Check-out: Belum\n\nKirim OTP Anda untuk *check-out*.", checkInTime)
	} else {
		checkInTime := utils.FormatTime(status.CheckInRecord.Timestamp, "HH:mm")
		checkOutTime := utils.FormatTime(status.CheckOutRecord.Timestamp, "HH:mm")
		duration := utils.CalculateWorkDuration(status.CheckInRecord.Timestamp, status.CheckOutRecord.Timestamp)
		message = fmt.Sprintf("✅ *Status Absensi*\n\n✅ Check-in: %s\n✅ Check-out: %s\n⌛ Durasi kerja: %s\n\nAbsensi hari ini sudah lengkap.", checkInTime, checkOutTime, duration)
	}

	return b.sendMarkdownMessage(msg.Chat.ID, message)
}

// handleAlias handles the /alias command
func (b *Bot) handleAlias(msg *Message, args []string) error {
	if len(args) == 0 {
		return b.sendMessage(msg.Chat.ID, "❌ Format tidak valid. Gunakan: /alias [Nama Depan] [Nama Belakang]")
	}

	firstName := utils.SanitizeName(args[0])
	if firstName == "" {
		return b.sendMessage(msg.Chat.ID, "❌ Nama depan tidak valid.")
	}

	var lastName *string
	if len(args) > 1 {
		lastNameVal := utils.SanitizeName(strings.Join(args[1:], " "))
		if lastNameVal != "" {
			lastName = &lastNameVal
		}
	}

	err := b.attendanceService.SetUserAlias(msg.From.ID, firstName, lastName)
	if err != nil {
		b.logger.Error("Failed to set user alias", "error", err, "user_id", msg.From.ID)
		return b.sendMessage(msg.Chat.ID, "❌ Gagal menyimpan alias. Silakan coba lagi.")
	}

	var aliasName string
	if lastName != nil {
		aliasName = fmt.Sprintf("%s %s", firstName, *lastName)
	} else {
		aliasName = firstName
	}

	return b.sendMessage(msg.Chat.ID, fmt.Sprintf("✅ Alias berhasil diatur: %s", aliasName))
}

// handleFullReport handles the /fullreport command
func (b *Bot) handleFullReport(msg *Message, args []string) error {
	response := `📊 *Laporan Lengkap Absensi*

Silakan masukkan password admin dan rentang tanggal dalam format:
` + "`[password] YYYY-MM-DD YYYY-MM-DD`" + `

*Contoh:*
` + "`admin123 2025-01-01 2025-01-31`" + `

*Catatan:* Laporan akan dikirim dalam format CSV.`

	// Set user session to await date range input
	b.sessions[msg.From.ID] = &SessionData{
		AwaitingDateRange: true,
	}

	return b.sendMarkdownMessage(msg.Chat.ID, response)
}

// handleOTP handles OTP verification and attendance marking
func (b *Bot) handleOTP(msg *Message) error {
	username := msg.From.Username
	if username == "" {
		username = fmt.Sprintf("user_%d", msg.From.ID)
	}

	firstName := utils.SanitizeName(msg.From.FirstName)
	var lastName *string
	if msg.From.LastName != "" {
		lastNameVal := utils.SanitizeName(msg.From.LastName)
		lastName = &lastNameVal
	}

	result, err := b.attendanceService.MarkAttendance(
		msg.From.ID,
		username,
		firstName,
		lastName,
		msg.Text,
	)
	if err != nil {
		b.logger.Error("Failed to mark attendance", "error", err, "user_id", msg.From.ID)
		return b.sendMessage(msg.Chat.ID, "❌ Terjadi kesalahan saat memproses absensi. Silakan coba lagi.")
	}

	if result.Success {
		return b.sendMarkdownMessage(msg.Chat.ID, result.Message)
	} else {
		return b.sendMessage(msg.Chat.ID, result.Message)
	}
}

// handleTextMessage handles non-command text messages
func (b *Bot) handleTextMessage(msg *Message) error {
	// Check if user is awaiting date range input for full report
	session := b.sessions[msg.From.ID]
	if session != nil && session.AwaitingDateRange {
		return b.handleFullReportInput(msg)
	}

	return b.sendMessage(msg.Chat.ID, "📝 Kirimkan kode OTP 6 digit Anda untuk absen, atau ketik /help untuk bantuan.")
}

// formatHistoryMessage formats attendance history into a readable message
func (b *Bot) formatHistoryMessage(records []models.AttendanceRecord) string {
	var message strings.Builder
	message.WriteString("📈 *Riwayat Absensi Anda (30 hari terakhir)*\n\n")

	// Group by date
	dailyRecords := make(map[string]map[string]*models.AttendanceRecord)
	dates := []string{}

	for _, record := range records {
		if dailyRecords[record.Date] == nil {
			dailyRecords[record.Date] = make(map[string]*models.AttendanceRecord)
			dates = append(dates, record.Date)
		}
		dailyRecords[record.Date][record.Type] = &record
	}

	// Sort dates in reverse order (newest first)
	for i := len(dates) - 1; i >= 0; i-- {
		date := dates[i]
		dayRecord := dailyRecords[date]

		// Parse and format date
		dateTime, err := utils.ParseDate(date)
		if err != nil {
			continue
		}
		displayDate := utils.FormatDate(dateTime, "dd MMMM yyyy")

		message.WriteString(fmt.Sprintf("%d. *%s*\n", len(dates)-i, displayDate))

		if checkIn := dayRecord["check_in"]; checkIn != nil {
			checkInTime := utils.FormatTime(checkIn.Timestamp, "HH:mm")
			status := " 🟢"
			if checkIn.Timestamp.Hour() >= 9 {
				status = " ⚠️"
			}
			message.WriteString(fmt.Sprintf("   ⏰ Masuk: %s%s\n", checkInTime, status))
		} else {
			message.WriteString("   ⏰ Masuk: -\n")
		}

		if checkOut := dayRecord["check_out"]; checkOut != nil {
			checkOutTime := utils.FormatTime(checkOut.Timestamp, "HH:mm")
			message.WriteString(fmt.Sprintf("   🏠 Pulang: %s\n", checkOutTime))
		} else {
			message.WriteString("   🏠 Pulang: -\n")
		}

		message.WriteString("\n")
	}

	uniqueDays := len(dates)
	totalRecords := len(records)

	message.WriteString("*Ringkasan:*\n")
	message.WriteString(fmt.Sprintf("📊 Total Hari: %d\n", uniqueDays))
	message.WriteString(fmt.Sprintf("📝 Total Absensi: %d", totalRecords))

	return message.String()
}

// handleFullReportInput processes user input for full report generation
func (b *Bot) handleFullReportInput(msg *Message) error {
	// Clear the session state
	delete(b.sessions, msg.From.ID)

	text := strings.TrimSpace(msg.Text)

	// Validate password and date range format
	dateRangeRegex := regexp.MustCompile(`^(\S+)\s+(\d{4}-\d{2}-\d{2})\s+(\d{4}-\d{2}-\d{2})$`)
	matches := dateRangeRegex.FindStringSubmatch(text)

	if len(matches) != 4 {
		return b.sendMessage(msg.Chat.ID, "❌ Format input tidak valid. Gunakan format: [password] YYYY-MM-DD YYYY-MM-DD\n\nContoh: admin123 2025-01-01 2025-01-31")
	}

	password := matches[1]
	startDate := matches[2]
	endDate := matches[3]

	// Check password
	if password != b.config.AdminPassword {
		return b.sendMessage(msg.Chat.ID, "❌ Password admin salah. Akses ditolak.")
	}

	// Validate dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return b.sendMessage(msg.Chat.ID, "❌ Tanggal mulai tidak valid. Pastikan format tanggal benar (YYYY-MM-DD).")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return b.sendMessage(msg.Chat.ID, "❌ Tanggal akhir tidak valid. Pastikan format tanggal benar (YYYY-MM-DD).")
	}

	if start.After(end) {
		return b.sendMessage(msg.Chat.ID, "❌ Tanggal mulai tidak boleh lebih besar dari tanggal akhir.")
	}

	// Generate and send CSV report
	if err := b.sendMessage(msg.Chat.ID, "⏳ Membuat laporan CSV... Mohon tunggu."); err != nil {
		return err
	}

	return b.generateAndSendCSVReport(msg.Chat.ID, startDate, endDate)
}

// generateAndSendCSVReport generates a CSV report and sends it as a document
func (b *Bot) generateAndSendCSVReport(chatID int64, startDate, endDate string) error {
	// Get attendance records for the date range
	records, err := b.attendanceService.GetAttendanceReportRange(startDate, endDate)
	if err != nil {
		b.logger.Error("Failed to get attendance records", "error", err)
		return b.sendMessage(chatID, "❌ Terjadi kesalahan saat mengambil data absensi.")
	}

	if len(records) == 0 {
		return b.sendMessage(chatID, "📭 Tidak ada data absensi dalam rentang tanggal yang ditentukan.")
	}

	// Generate CSV file
	filePath, err := b.csvGenerator.GenerateAttendanceReport(records, startDate, endDate)
	if err != nil {
		b.logger.Error("Failed to generate CSV report", "error", err)
		return b.sendMessage(chatID, "❌ Terjadi kesalahan saat membuat laporan CSV.")
	}

	// Send CSV file
	file, err := os.Open(filePath)
	if err != nil {
		b.logger.Error("Failed to open CSV file", "error", err)
		return b.sendMessage(chatID, "❌ Terjadi kesalahan saat membuka file laporan.")
	}
	defer file.Close()

	filename := fmt.Sprintf("attendance_%s_to_%s.csv", startDate, endDate)

	// Send the file
	if err := b.api.SendDocument(chatID, file, filename); err != nil {
		b.logger.Error("Failed to send CSV document", "error", err)
		return b.sendMessage(chatID, "❌ Terjadi kesalahan saat mengirim laporan.")
	}

	// Send confirmation message with statistics
	caption := fmt.Sprintf("📊 *Laporan Absensi*\n\n📅 Periode: %s s/d %s\n📈 Total Records: %d",
		startDate, endDate, len(records))

	// Clean up temp file
	if err := os.Remove(filePath); err != nil {
		b.logger.Warn("Failed to clean up temp file", "file", filePath, "error", err)
	}

	return b.sendMarkdownMessage(chatID, caption)
}

// sendMessage sends a plain text message
func (b *Bot) sendMessage(chatID int64, text string) error {
	return b.api.SendMessage(chatID, text)
}

// sendMarkdownMessage sends a message with Markdown formatting
func (b *Bot) sendMarkdownMessage(chatID int64, text string) error {
	options := &SendMessageOptions{
		ParseMode: "Markdown",
	}
	return b.api.SendMessageWithOptions(chatID, text, options)
}
