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
			return newValidationError("file does not exist in directory")
		}
		return newValidationError("file access error")
	}
	return nil
}

func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("failed to read file")
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
				return "", newValidationError("invalid character in line")
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", newValidationError("line must be empty")
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
			return "", newValidationError("line cannot be empty")
		}
		if len(trimmed) != 4 {
			return "", newValidationError("line must have exactly 4 characters")
		}

		if blockIndex >= 4 {
			return "", newValidationError("block has too many lines")
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
		return "", newValidationError("file is empty")
	}
	if lineCount < minLines {
		return "", newValidationError("file must have at least lines")
	}

	return SolveTetrominos(tetrominos)
}