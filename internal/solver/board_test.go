package solver

import (
	"testing"
)

func makeTetromino(letter rune, points []Point) *Tetromino {
	return &Tetromino{
		Letter: letter,
		Points: points,
	}
}

func TestNewBoard(t *testing.T) {
	b := NewBoard(4)
	if b == nil {
		t.Fatal("Expected non-nil board")
	}
	if b.Size != 4 {
		t.Errorf("Expected board size 4, got %d", b.Size)
	}
	if len(b.Grid) != 4 || len(b.Grid[0]) != 4 {
		t.Error("Board grid not properly initialized")
	}

	nilBoard := NewBoard(0)
	if nilBoard != nil {
		t.Error("Expected nil board for size 0")
	}
}

func TestCanPlace(t *testing.T) {
	board := NewBoard(4)
	tetromino := makeTetromino('A', []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}) // 2x2 block

	// Valid placement
	if !board.CanPlace(tetromino, 1, 1) {
		t.Error("Expected valid placement at (1,1)")
	}

	// Out of bounds placement
	if board.CanPlace(tetromino, 3, 3) {
		t.Error("Expected invalid placement (out of bounds)")
	}

	// Overlap placement
	board.Place(tetromino, 1, 1)
	if board.CanPlace(tetromino, 1, 1) {
		t.Error("Expected invalid placement (overlap)")
	}
}
