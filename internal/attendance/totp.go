package attendance

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// TOTPService handles Time-based One-Time Password operations
type TOTPService struct {
	secret string
}

// NewTOTPService creates a new TOTP service with the given secret
func NewTOTPService(secret string) *TOTPService {
	return &TOTPService{
		secret: secret,
	}
}

// Verify checks if the provided token is valid for the current time
func (t *TOTPService) Verify(token string) bool {
	// Remove any spaces or formatting
	token = strings.ReplaceAll(token, " ", "")

	if len(token) != 6 {
		return false
	}

	// Check current time and ±1 time step for clock skew tolerance
	now := time.Now().Unix()
	timeStep := int64(30) // 30 seconds

	for i := -1; i <= 1; i++ {
		testTime := (now/timeStep + int64(i)) * timeStep
		expectedToken := t.generateTOTPForTime(testTime)
		if token == expectedToken {
			return true
		}
	}

	return false
}

// Generate creates a TOTP token for the current time
func (t *TOTPService) Generate() string {
	now := time.Now().Unix()
	return t.generateTOTPForTime(now)
}

// generateTOTPForTime creates a TOTP token for a specific time
func (t *TOTPService) generateTOTPForTime(unixTime int64) string {
	timeStep := int64(30) // 30 seconds
	counter := unixTime / timeStep

	// Convert secret from base32
	secret, err := base32.StdEncoding.DecodeString(strings.ToUpper(t.secret))
	if err != nil {
		return ""
	}

	// Create HMAC-SHA1 hash
	h := hmac.New(sha1.New, secret)

	// Convert counter to bytes
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	h.Write(counterBytes)
	hash := h.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0f

	// Extract 4 bytes starting from offset
	truncatedHash := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	// Generate 6-digit code
	code := truncatedHash % 1000000

	return fmt.Sprintf("%06d", code)
}

// GenerateSecret creates a new random base32-encoded secret
func GenerateSecret() string {
	// Generate 20 random bytes (160 bits)
	secretBytes := make([]byte, 20)
	for i := range secretBytes {
		secretBytes[i] = byte(rand.Intn(256))
	}

	// Encode as base32
	return base32.StdEncoding.EncodeToString(secretBytes)
}

// GenerateKeyURI creates an otpauth:// URI for use with authenticator apps
func (t *TOTPService) GenerateKeyURI(accountName, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		issuer, accountName, t.secret, issuer)
}

// GetTimeRemaining returns the number of seconds until the current TOTP expires
func (t *TOTPService) GetTimeRemaining() int {
	now := time.Now().Unix()
	timeStep := int64(30)
	return int(timeStep - (now % timeStep))
}

// ValidateSecret checks if a secret is properly formatted
func ValidateSecret(secret string) bool {
	if len(secret) < 16 {
		return false
	}

	// Try to decode as base32
	_, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	return err == nil
}
