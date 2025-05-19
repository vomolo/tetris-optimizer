package solver

import (
<<<<<<< HEAD
	"bytes"
=======
>>>>>>> 6a47205 (REFINE: Simplify validation error handling and improve comments in validator)
	"fmt"
	"os"
	"path/filepath" // helps with file path manipulation for cross-platform compatibility
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
<<<<<<< HEAD
<<<<<<< HEAD

func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("ERROR")
	}

	if len(content) < 16 {
		return "", newValidationError("ERROR")
	}

	lines := bytes.Split(content, []byte{'\n'})
	var (
		lineCount    int
		blockLines   [4][]byte
=======

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
>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
		blockIndex   int
		hasContent   bool
		blockCounter int
		tetrominos   []*Tetromino
	)

	for _, line := range lines {
		lineCount++
<<<<<<< HEAD
		trimmed := bytes.TrimSpace(line)

		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", newValidationError("ERROR")
=======
		trimmed := strings.TrimSpace(line)

		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", NewValidationError("invalid character in input")
>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
<<<<<<< HEAD
				return "", newValidationError("ERROR")
			}

			tetromino, err := validateAndCreateTetromino(blockLines[:], blockCounter)
			if err != nil {
				return "", err
			}
			tetrominos = append(tetrominos, tetromino)
			blockIndex = 0
			blockCounter++
			continue
		}

		if len(trimmed) == 0 {
			return "", newValidationError("ERROR")
		}
		if len(trimmed) != 4 {
			return "", newValidationError("ERROR")
		}

		if blockIndex >= 4 {
			return "", newValidationError("ERROR")
=======
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
>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
		}
		blockLines[blockIndex] = trimmed
		blockIndex++
		hasContent = true
	}

	if blockIndex > 0 {
<<<<<<< HEAD
		tetromino, err := validateAndCreateTetromino(blockLines[:blockIndex], blockCounter)
=======
		tetromino, err := validateAndCreateTetrominoStr(blockLines[:blockIndex], blockCounter)
>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
		if err != nil {
			return "", err
		}
		tetrominos = append(tetrominos, tetromino)
	}

<<<<<<< HEAD
	if !hasContent {
		return "", newValidationError("ERROR")
	}
	if lineCount < minLines {
		return "", newValidationError("ERROR")
=======
	if !hasContent || lineCount < minLines {
		return "", NewValidationError("input lacks valid tetrominos")
>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
	}

	return SolveTetrominos(tetrominos)
}
<<<<<<< HEAD
=======
>>>>>>> 6a47205 (REFINE: Simplify validation error handling and improve comments in validator)
=======

>>>>>>> 0e51e2e (REFINE: Enhance validation logic for tetromino content and improve error handling)
