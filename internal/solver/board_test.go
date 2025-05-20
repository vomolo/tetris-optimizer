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

func TestPlaceAndRemove(t *testing.T) {
	board := NewBoard(4)
	tetromino := makeTetromino('B', []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}})

	board.Place(tetromino, 0, 0)
	if board.Grid[0][0] != 'B' || board.Grid[1][1] != 'B' {
		t.Error("Tetromino not placed correctly")
	}
	if board.Placed != 1 {
		t.Errorf("Expected Placed = 1, got %d", board.Placed)
	}

	board.Remove(tetromino, 0, 0)
	if board.Grid[0][0] != 0 || board.Grid[1][1] != 0 {
		t.Error("Tetromino not removed correctly")
	}
	if board.Placed != 0 {
		t.Errorf("Expected Placed = 0, got %d", board.Placed)
	}
}

func TestString(t *testing.T) {
	board := NewBoard(2)
	tetromino := makeTetromino('C', []Point{{0, 0}})
	board.Place(tetromino, 0, 0)

	expected := "C.\n.."
	got := board.String()

	if got != expected {
		t.Errorf("Expected board string:\n%s\nGot:\n%s", expected, got)
	}

	board.Remove(tetromino, 0, 0)
	expectedEmpty := "..\n.."
	if board.String() != expectedEmpty {
		t.Errorf("Expected empty board:\n%s\nGot:\n%s", expectedEmpty, board.String())
	}
}
