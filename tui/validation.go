package tui

import (
	"fmt"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

// FilenameValidator validates filenames for export operations.
//
// This validator ensures filenames are safe for use in filesystem operations,
// preventing path traversal, shell injection, and other security issues.
//
// Example usage:
//
//	validator := NewFilenameValidator()
//	if err := validator.Validate("my-file.txt"); err != nil {
//	    log.Fatal(err)
//	}
//
//	safeName := validator.SanitizeFilename("unsafe/../file.txt")
//	// safeName will be "file.txt" with path components removed
type FilenameValidator struct {
	// MaxLength is the maximum allowed filename length in bytes (default: 255)
	MaxLength int
	// AllowedPattern is the regex pattern for valid characters
	AllowedPattern *regexp.Regexp
}

// NewFilenameValidator creates a validator with secure default settings.
//
// Default configuration:
//   - MaxLength: 255 bytes (typical filesystem limit)
//   - AllowedPattern: ^[a-zA-Z0-9_. -]+$ (alphanumeric, underscore, dot, space, hyphen)
//
// Returns:
//   - A configured FilenameValidator ready for use
func NewFilenameValidator() *FilenameValidator {
	return &FilenameValidator{
		MaxLength:      255, // Typical filesystem limit
		AllowedPattern: regexp.MustCompile(`^[a-zA-Z0-9_. -]+$`),
	}
}

// Validate checks if a filename is safe and valid.
//
// The validation performs multiple checks:
//   - UTF-8 validity
//   - Length limits
//   - Special directory names (".", "..")
//   - Character whitelist (alphanumeric, underscore, dot, space, hyphen)
//
// Parameters:
//   - filename: The filename to validate
//
// Returns:
//   - nil if valid
//   - error describing the validation failure
//
// Example:
//
//	err := validator.Validate("../../../etc/passwd")
//	// Returns error: "filename contains unsafe characters"
func (v *FilenameValidator) Validate(filename string) error {
	// Check UTF-8 validity
	if !utf8.ValidString(filename) {
		return fmt.Errorf("filename contains invalid UTF-8")
	}

	// Check length
	if len(filename) > v.MaxLength {
		return fmt.Errorf("filename too long: %d bytes (max %d)", len(filename), v.MaxLength)
	}

	if len(filename) == 0 {
		return fmt.Errorf("filename is empty")
	}

	// Check for special names
	if filename == "." || filename == ".." {
		return fmt.Errorf("filename cannot be . or ..")
	}

	// Check allowed characters
	if !v.AllowedPattern.MatchString(filename) {
		return fmt.Errorf("filename contains unsafe characters")
	}

	return nil
}

// SanitizeFilename attempts to make a filename safe by removing or replacing unsafe characters.
//
// The sanitization process:
//  1. Strips directory components using filepath.Base()
//  2. Replaces unsafe characters with underscores
//  3. Truncates to maximum length if needed
//
// Only alphanumeric characters, spaces, hyphens, underscores, and dots are preserved.
// All other characters are replaced with underscores.
//
// Parameters:
//   - filename: The potentially unsafe filename
//
// Returns:
//   - A sanitized filename safe for use
//
// Example:
//
//	safe := validator.SanitizeFilename("../bad;file|name$.txt")
//	// Returns: "bad_file_name_.txt"
//
//	safe := validator.SanitizeFilename("/path/to/file.txt")
//	// Returns: "file.txt"
func (v *FilenameValidator) SanitizeFilename(filename string) string {
	// Remove any directory components
	filename = filepath.Base(filename)

	// Replace unsafe characters with underscores
	safe := make([]rune, 0, len(filename))
	for _, r := range filename {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			safe = append(safe, r)
		case r == ' ' || r == '-' || r == '_' || r == '.':
			safe = append(safe, r)
		default:
			safe = append(safe, '_')
		}
	}

	result := string(safe)

	// Truncate if too long
	if len(result) > v.MaxLength {
		result = result[:v.MaxLength]
	}

	return result
}

// ValidateAndSanitize performs validation and returns a sanitized version if validation fails.
//
// This is a convenience method that combines Validate and SanitizeFilename.
// If the filename is already valid, it returns the original.
// If invalid, it returns a sanitized version.
//
// Parameters:
//   - filename: The filename to validate and potentially sanitize
//
// Returns:
//   - sanitized: A safe filename (either original or sanitized)
//   - modified: true if sanitization was needed, false if original was valid
//
// Example:
//
//	safe, modified := validator.ValidateAndSanitize("good-file.txt")
//	// Returns: "good-file.txt", false
//
//	safe, modified := validator.ValidateAndSanitize("bad;file.txt")
//	// Returns: "bad_file.txt", true
func (v *FilenameValidator) ValidateAndSanitize(filename string) (sanitized string, modified bool) {
	if err := v.Validate(filename); err == nil {
		return filename, false
	}
	return v.SanitizeFilename(filename), true
}
