package functions

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	tetrisDir = "tetris_files"
	minLines  = 4
)

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func newValidationError(format string, args ...interface{}) error {
	return &ValidationError{
		message: fmt.Sprintf(format, args...),
	}
}

// Validate checks the tetromino file and returns formatted output with letters
func Validate(filename string) (string, error) {
	fullPath := filename
	if filepath.Dir(filename) != tetrisDir {
		fullPath = filepath.Join(tetrisDir, filename)
	}

	if err := validateStructure(fullPath); err != nil {
		return "", err
	}
	return validateAndPrintContent(fullPath)
}

func validateStructure(fullPath string) error {
	if filepath.Ext(fullPath) != ".txt" {
		return newValidationError("file must have .txt extension")
	}

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return newValidationError("file '%s' does not exist in %s directory",
				filepath.Base(fullPath), tetrisDir)
		}
		return newValidationError("file access error: %v", err)
	}
	return nil
}

func validateAndPrintContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("failed to read file: %v", err)
	}

	if len(content) < 20 {
		return "", newValidationError("file too small to be valid")
	}

	lines := bytes.Split(content, []byte{'\n'})
	var (
		lineCount    int
		blockLines   [4][]byte
		blockIndex   int
		hasContent   bool
		blockCounter int
		output       bytes.Buffer
	)

	for _, line := range lines {
		lineCount++
		trimmed := bytes.TrimSpace(line)

		// Validate characters
		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", newValidationError("invalid character '%c' in line %d", char, lineCount)
			}
		}

		// Handle separator lines
		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", newValidationError("line %d must be empty", lineCount)
			}

			if err := processBlock(blockLines[:], blockCounter, &output); err != nil {
				return "", err
			}
			blockIndex = 0
			blockCounter++
			continue
		}

		// Validate tetromino lines
		if len(trimmed) == 0 {
			return "", newValidationError("line %d cannot be empty", lineCount)
		}
		if len(trimmed) != 4 {
			return "", newValidationError("line %d must have exactly 4 characters", lineCount)
		}

		if blockIndex >= 4 {
			return "", newValidationError("block %d has too many lines", blockCounter+1)
		}
		blockLines[blockIndex] = trimmed
		blockIndex++
		hasContent = true
	}

	// Check last block if no trailing separator
	if blockIndex > 0 {
		if err := processBlock(blockLines[:blockIndex], blockCounter, &output); err != nil {
			return "", err
		}
	}

	if !hasContent {
		return "", newValidationError("file is empty")
	}
	if lineCount < minLines {
		return "", newValidationError("file must have at least %d lines", minLines)
	}
	return output.String(), nil
}

func processBlock(block [][]byte, blockNumber int, output *bytes.Buffer) error {
	if len(block) != 4 {
		return newValidationError("block must have 4 lines (found %d)", len(block))
	}

	var (
		hashCount int
		shape     uint16
	)

	// Count hashes and build shape
	for y, line := range block {
		for x, char := range line {
			if char == '#' {
				hashCount++
				shape |= 1 << uint(y*4+x)
			}
		}
	}

	if hashCount != 4 {
		return newValidationError("block must have exactly 4 '#' (found %d)", hashCount)
	}

	// Check if valid tetromino
	validShapes := map[uint16]bool{
		0xF000: true, 0x1111: true, 0xCC00: true,
		0xE400: true, 0x4C40: true, 0x4E00: true, 0x8C80: true,
		0xE800: true, 0xC440: true, 0x2E00: true, 0x88C0: true,
		0xE200: true, 0x44C0: true, 0x8E00: true, 0xC880: true,
		0x6C00: true, 0x8C40: true, 0xC600: true, 0x4C80: true,
	}

	if !validShapes[shape] {
		return newValidationError("block is not a valid tetromino")
	}

	// Replace # with corresponding letter and print
	letter := 'A' + rune(blockNumber)
	for _, line := range block {
		for _, char := range line {
			if char == '#' {
				output.WriteRune(letter)
			} else {
				output.WriteByte(char)
			}
		}
		output.WriteByte('\n')
	}
	output.WriteByte('\n') // Add separator after each block

	return nil
}
