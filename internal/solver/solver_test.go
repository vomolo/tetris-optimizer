package solver

import (
	"testing"
)

// TestNewBoard tests board creation.
func TestNewBoard(t *testing.T) {
	tests := []struct {
		size      int
		shouldNil bool
	}{
		{0, true},
		{-1, true},
		{4, false},
	}

	for _, tt := range tests {
		b := NewBoard(tt.size)
		if tt.shouldNil && b != nil {
			t.Errorf("NewBoard(%d) should return nil", tt.size)
		}
		if !tt.shouldNil && (b == nil || b.Size != tt.size || len(b.Grid) != tt.size) {
			t.Errorf("NewBoard(%d) created invalid board", tt.size)
		}
	}
}

// TestBoardPlacement tests placing and removing tetrominos.
func TestBoardPlacement(t *testing.T) {
	b := NewBoard(4)
	tet := &Tetromino{
		Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		Letter: 'A',
		Width:  2,
		Height: 2,
	}

	// Test valid placement
	if !b.CanPlace(tet, 0, 0) {
		t.Error("CanPlace should allow placement at (0,0)")
	}
	b.Place(tet, 0, 0)
	if b.Placed != 1 || b.Grid[0][0] != 'A' || b.Grid[1][1] != 'A' {
		t.Error("Place failed to update board correctly")
	}

	// Test invalid placement
	if b.CanPlace(tet, 0, 0) {
		t.Error("CanPlace should prevent overlap")
	}

	// Test removal
	b.Remove(tet, 0, 0)
	if b.Placed != 0 || b.Grid[0][0] != 0 {
		t.Error("Remove failed to clear board")
	}
}
