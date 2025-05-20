package solver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tetrisDir = "testfiles"
	minLines  = 4
)

// Validate validates a Tetris input file and returns the solved board.
func Validate(filename string) (string, error) {
	// Clean the input path
	cleanFilename := filepath.Clean(filename)

	// Get absolute path of the tetris directory
	absTetrisDir, err := filepath.Abs(tetrisDir)
	if err != nil {
		return "", fmt.Errorf("invalid tetris directory: %v", err)
	}

	// Get absolute path of the requested file
	absFilePath, err := filepath.Abs(cleanFilename)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %v", err)
	}

	// If file isn't already in the tetris directory, join them
	if !strings.HasPrefix(absFilePath, absTetrisDir+string(filepath.Separator)) {
		absFilePath = filepath.Join(absTetrisDir, cleanFilename)
	}

	// Prevent directory traversal
	if !strings.HasPrefix(absFilePath, absTetrisDir+string(filepath.Separator)) {
		return "", NewValidationError("invalid file path: attempted directory traversal")
	}

	if err := validateStructure(absFilePath); err != nil {
		return "", err
	}
	return validateAndSolveContent(absFilePath)
}

// validateStructure checks the file's structure (extension and existence).
func validateStructure(fullPath string) error {
	if filepath.Ext(fullPath) != ".txt" {
		return NewValidationError("file must have .txt extension")
	}

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return NewValidationError("file does not exist in directory")
		}
		return NewValidationError("file access error")
	}
	return nil
}

// validateAndSolveContent reads and processes the file content.
func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", NewValidationError("error reading file")
	}
	return validateAndSolve(string(content))
}

// validateAndSolve validates the content and solves the tetromino puzzle.
func validateAndSolve(content string) (string, error) {
	if len(content) < 16 {
		return "", NewValidationError("content too short")
	}

	lines := strings.Split(content, "\n")
	var (
		lineCount    int
		blockLines   [4]string
		blockIndex   int
		hasContent   bool
		blockCounter int
		tetrominos   []*Tetromino
	)

	for _, line := range lines {
		lineCount++
		trimmed := strings.TrimSpace(line)

		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", NewValidationError("invalid character in input")
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", NewValidationError("separator line must be empty")
			}
			if blockIndex == 4 {
				tetromino, err := validateAndCreateTetrominoStr(blockLines[:], blockCounter)
				if err != nil {
					return "", err
				}
				tetrominos = append(tetrominos, tetromino)
				blockCounter++
			}
			blockIndex = 0
			continue
		}

		if len(trimmed) == 0 || len(trimmed) != 4 {
			return "", NewValidationError("each tetromino line must be 4 characters")
		}

		if blockIndex >= 4 {
			return "", NewValidationError("tetromino has too many lines")
		}
		blockLines[blockIndex] = trimmed
		blockIndex++
		hasContent = true
	}

	if blockIndex > 0 {
		tetromino, err := validateAndCreateTetrominoStr(blockLines[:blockIndex], blockCounter)
		if err != nil {
			return "", err
		}
		tetrominos = append(tetrominos, tetromino)
	}

	if !hasContent || lineCount < minLines {
		return "", NewValidationError("input lacks valid tetrominos")
	}

	return SolveTetrominos(tetrominos)
}

// validateAndCreateTetrominoStr converts string lines to a tetromino.
func validateAndCreateTetrominoStr(lines []string, id int) (*Tetromino, error) {
	byteLines := make([][]byte, len(lines))
	for i, line := range lines {
		byteLines[i] = []byte(line)
	}
	return ValidateAndCreateTetromino(byteLines, id)
}