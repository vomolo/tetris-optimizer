package functions

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	tetrisDir = "tetris_files"  // The required directory
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

// validateContent checks the file content rules (unchanged)
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

		if lineCount%5 == 0 && len(bytes.TrimSpace(line)) > 0 {
			return newValidationError("line %d must be empty", lineCount)
		}

		if !hasContent && len(bytes.TrimSpace(line)) > 0 {
			hasContent = true
		}

		if lineCount >= minLines && hasContent {
			break
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