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
		return newValidationError("File must have .txt extension")
	}

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return newValidationError("File does not exist in directory")
		}
		return newValidationError("File access error")
	}
	return nil
}

func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("Failed to read file")
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
				return "", newValidationError("ERROR")
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", newValidationError("ERROR")
			}

			tetromino := validateAndCreateTetromino(blockLines[:], blockCounter)
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
		}
		blockLines[blockIndex] = trimmed
		blockIndex++
		hasContent = true
	}

	if blockIndex > 0 {
		tetromino := validateAndCreateTetromino(blockLines[:blockIndex], blockCounter)
		if err != nil {
			return "", err
		}
		tetrominos = append(tetrominos, tetromino)
	}

	if !hasContent {
		return "", newValidationError("ERROR")
	}
	if lineCount < minLines {
		return "", newValidationError("ERROR")
	}

	return SolveTetrominos(tetrominos)
}