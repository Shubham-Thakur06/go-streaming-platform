package validator

import (
	"regexp"
	"strings"
)

// Validator provides validation functionality
type Validator struct{}

// New creates a new validator instance
func New() *Validator {
	return &Validator{}
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername validates a username
func (v *Validator) ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword validates a password
func (v *Validator) ValidatePassword(password string) bool {
	return len(password) >= 6
}

// ValidateFileType validates if a file type is allowed
func (v *Validator) ValidateFileType(filename string) bool {
	allowedTypes := []string{".mp3", ".mp4", ".wav", ".avi", ".mov", ".mkv"}
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])

	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// SanitizeString sanitizes a string input
func (v *Validator) SanitizeString(input string) string {
	return strings.TrimSpace(input)
}
