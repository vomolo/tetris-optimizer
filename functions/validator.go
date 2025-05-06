package functions

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	tetrisDir = "tetris_files" // The required directory
	minLines  = 4
)

// Validate automatically prepends tetrisDir to the filename
func Validate(filename string) error {
	// Prepend tetrisDir if not already present
	fullPath := filename
	if filepath.Dir(filename) != tetrisDir {
		fullPath = filepath.Join(tetrisDir, filename)
	}

	if err := validateStructure(fullPath); err != nil {
		return err
	}
	return validateContent(fullPath)
}

// validateStructure checks the file exists and has .txt extension
func validateStructure(fullPath string) error {
	// Check extension first (fastest check)
	if filepath.Ext(fullPath) != ".txt" {
		return newValidationError("file must have .txt extension")
	}

	// Absolute path check
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return newValidationError("path resolution failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return newValidationError("file '%s' does not exist in %s directory",
				filepath.Base(fullPath), tetrisDir)
		}
		return newValidationError("file access error: %v", err)
	}

	return nil
}

// validateContent checks the file content rules with strict character and length validation
func validateContent(fullPath string) error {
	file, err := os.Open(fullPath)
	if err != nil {
		return newValidationError("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024) // 64KB buffer
	scanner.Buffer(buf, cap(buf))

	var (
		lineCount  int
		hasContent bool
	)

	for scanner.Scan() {
		lineCount++
		line := scanner.Bytes()
		trimmedLine := bytes.TrimSpace(line)

		// Validate all characters in the line
		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return newValidationError("invalid character '%c' in line %d - only '#' and '.' are allowed", char, lineCount)
			}
		}

		// Every 5th line must be empty (separator between tetrominoes)
		if lineCount%5 == 0 {
			if len(trimmedLine) > 0 {
				return newValidationError("line %d must be empty (separator line)", lineCount)
			}
			continue // Skip further checks for separator lines
		}

		// For non-separator lines
		if len(trimmedLine) == 0 {
			return newValidationError("line %d cannot be empty (contains tetromino data)", lineCount)
		}

		// Strict 4-character check for non-empty lines
		if len(trimmedLine) != 4 {
			return newValidationError("line %d must have exactly 4 characters (found %d: %q)",
				lineCount, len(trimmedLine), trimmedLine)
		}

		// Mark that we've found content
		if !hasContent {
			hasContent = true
		}
	}

	if err := scanner.Err(); err != nil {
		return newValidationError("file read error: %v", err)
	}

	if !hasContent {
		return newValidationError("file is empty")
	}

	if lineCount < minLines {
		return newValidationError("file must have at least %d lines", minLines)
	}

	return nil
}

func newValidationError(format string, args ...interface{}) error {
	return &ValidationError{
		message: fmt.Sprintf(format, args...),
	}
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}
