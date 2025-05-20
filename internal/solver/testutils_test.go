package solver

import (
	"reflect"
	"strings"
)

// createTestTetromino creates a tetromino from a 4x4 string representation.
func createTestTetromino(shape []string, id int) (*Tetromino, error) {
	if len(shape) != 4 {
		return nil, newValidationError("test tetromino must have 4 rows")
	}
	byteLines := make([][]byte, 4)
	for i, line := range shape {
		if len(line) != 4 {
			return nil, newValidationError("test tetromino row must have 4 columns")
		}
		byteLines[i] = []byte(line)
	}
	return ValidateAndCreateTetromino(byteLines, id)
}

// compareBoards compares two board string representations, ignoring trailing newlines.
func compareBoards(actual, expected string) bool {
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")
	return reflect.DeepEqual(actualLines, expectedLines)
}