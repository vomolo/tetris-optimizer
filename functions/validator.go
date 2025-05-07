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
		lineCount    int
		blockLines   []string
		hasContent   bool
		blockCounter int
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

			// Validate the completed block
			if err := validateTetrominoBlock(blockLines, blockCounter+1); err != nil {
				return err
			}
			blockLines = blockLines[:0] // Reset block
			blockCounter++
			continue
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

		// Add to current block
		blockLines = append(blockLines, string(trimmedLine))
		hasContent = true
	}

	if err := scanner.Err(); err != nil {
		return newValidationError("file read error: %v", err)
	}

	// Check the last block if there was no trailing separator
	if len(blockLines) > 0 {
		if err := validateTetrominoBlock(blockLines, blockCounter+1); err != nil {
			return err
		}
	}

	if !hasContent {
		return newValidationError("file is empty")
	}

	if lineCount < minLines {
		return newValidationError("file must have at least %d lines", minLines)
	}

	return nil
}

// validateTetrominoBlock checks if a 4-line block contains exactly one valid tetromino
func validateTetrominoBlock(block []string, blockNumber int) error {
	if len(block) != 4 {
		return newValidationError("block %d must have exactly 4 lines (found %d)", blockNumber, len(block))
	}

	// Count the number of '#' characters
	var hashCount int
	for _, line := range block {
		for _, char := range line {
			if char == '#' {
				hashCount++
			}
		}
	}

	if hashCount != 4 {
		return newValidationError("block %d must contain exactly 4 '#' characters (found %d)", blockNumber, hashCount)
	}

	// Convert block to grid for adjacency check
	grid := make([][]rune, 4)
	for i, line := range block {
		grid[i] = []rune(line)
	}

	// Check if the '#' characters form a single connected tetromino
	if !isValidTetromino(grid) {
		return newValidationError("block %d does not form a valid tetromino", blockNumber)
	}

	return nil
}

// isValidTetromino checks if the grid contains exactly one valid tetromino
func isValidTetromino(grid [][]rune) bool {
	// Find first '#' to start DFS
	var startX, startY int
	found := false
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if grid[y][x] == '#' {
				startX, startY = x, y
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return false
	}

	// Perform DFS to count connected '#' characters
	visited := make([][]bool, 4)
	for i := range visited {
		visited[i] = make([]bool, 4)
	}

	count := 0
	var dfs func(x, y int)
	dfs = func(x, y int) {
		if x < 0 || x >= 4 || y < 0 || y >= 4 || visited[y][x] || grid[y][x] != '#' {
			return
		}
		visited[y][x] = true
		count++
		dfs(x+1, y)
		dfs(x-1, y)
		dfs(x, y+1)
		dfs(x, y-1)
	}

	dfs(startX, startY)

	// All '#' should be connected
	return count == 4
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
