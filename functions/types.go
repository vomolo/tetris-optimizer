package functions

import "fmt"

type Point struct {
	X, Y int
}

type Tetromino struct {
	Points []Point
	Letter rune
	Width  int
	Height int
}

type Board struct {
	Grid   [][]rune
	Size   int
	Placed int
}

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
