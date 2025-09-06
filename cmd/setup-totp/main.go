package main

import (
	"attendance-bot/internal/attendance"
	"fmt"
	"os"
)

func main() {
	// Generate a new TOTP secret
	secret := attendance.GenerateSecret()
	fmt.Printf("Generated TOTP Secret: %s\n", secret)

	// Create TOTP service
	totpService := attendance.NewTOTPService(secret)

	// Service details
	serviceName := "Attendance Bot"
	accountName := "Employee"

	// Generate otpauth URL
	otpauthURL := totpService.GenerateKeyURI(accountName, serviceName)
	fmt.Printf("OTP Auth URL: %s\n", otpauthURL)

	// Print QR code information
	fmt.Println("\n=== Setup Instructions ===")
	fmt.Printf("1. Copy the TOTP Secret above to your .env file as TOTP_SECRET=%s\n", secret)
	fmt.Println("2. Use a QR code generator to create a QR code with the OTP Auth URL above")
	fmt.Println("3. Scan the QR code with your authenticator app (Google Authenticator, Authy, etc.)")
	fmt.Println("4. Or manually enter the secret in your authenticator app")
	fmt.Println("5. Start the bot and test with the 6-digit code from your app")

	// Generate current TOTP token for testing
	currentToken := totpService.Generate()
	fmt.Printf("\nCurrent TOTP token (for testing): %s\n", currentToken)

	fmt.Printf("Time remaining for current token: %d seconds\n", totpService.GetTimeRemaining())

	// Create a sample .env file content
	fmt.Println("\n=== Sample .env file ===")
	fmt.Printf("BOT_TOKEN=your_telegram_bot_token_here\n")
	fmt.Printf("TOTP_SECRET=%s\n", secret)
	fmt.Printf("ADMIN_PASSWORD=your_admin_password_here\n")
	fmt.Printf("NODE_ENV=development\n")

	// Write to .env.example if it doesn't exist
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		envContent := fmt.Sprintf(`# Telegram Bot Configuration
BOT_TOKEN=your_telegram_bot_token_here

# TOTP Secret for attendance verification
TOTP_SECRET=%s

# Admin password for management functions
ADMIN_PASSWORD=your_admin_password_here

# Environment (development or production)
NODE_ENV=development

# Database path (optional, defaults to data/attendance.db)
DATABASE_PATH=data/attendance.db
`, secret)

		if err := os.WriteFile(".env.example", []byte(envContent), 0644); err == nil {
			fmt.Println("\n✅ Created .env.example file with generated secret")
		}
	}
}
