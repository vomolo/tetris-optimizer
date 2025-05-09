package solver

import (
	"bytes"
	"os"
	"path/filepath"
)

const (
	tetrisDir = "testfiles"
	minLines  = 4
)

func Validate(filename string) (string, error) {
	fullPath := filename
	if filepath.Dir(filename) != tetrisDir {
		fullPath = filepath.Join(tetrisDir, filename)
	}

	if err := validateStructure(fullPath); err != nil {
		return "", err
	}
	return validateAndSolveContent(fullPath)
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

func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("failed to read file: %v", err)
	}

	if len(content) < 16 {
		return "", newValidationError("ERROR")
	}

	lines := bytes.Split(content, []byte{'\n'})
	var (
		lineCount    int
		blockLines   [4][]byte
		blockIndex   int
		hasContent   bool
		blockCounter int
		tetrominos   []*Tetromino
	)

	for _, line := range lines {
		lineCount++
		trimmed := bytes.TrimSpace(line)

		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", newValidationError("invalid character '%c' in line %d", char, lineCount)
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", newValidationError("line %d must be empty", lineCount)
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

	if blockIndex > 0 {
		tetromino, err := validateAndCreateTetromino(blockLines[:blockIndex], blockCounter)
		if err != nil {
			return "", err
		}
		tetrominos = append(tetrominos, tetromino)
	}

	if !hasContent {
		return "", newValidationError("file is empty")
	}
	if lineCount < minLines {
		return "", newValidationError("file must have at least %d lines", minLines)
	}

	return SolveTetrominos(tetrominos)
}