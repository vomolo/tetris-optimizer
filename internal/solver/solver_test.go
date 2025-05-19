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

// TestValidateAndCreateTetromino tests tetromino creation.
func TestValidateAndCreateTetromino(t *testing.T) {
	tests := []struct {
		block     [][]byte
		wantError bool
	}{
		{ // Valid square tetromino
			[][]byte{
				[]byte("##.."),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			false,
		},
		{ // Invalid: too many blocks
			[][]byte{
				[]byte("###."),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			true,
		},
		{ // Invalid: not connected
			[][]byte{
				[]byte("#..."),
				[]byte("...."),
				[]byte("...."),
				[]byte("#..#"),
			},
			true,
		},
		{ // Invalid: wrong size
			[][]byte{
				[]byte("##.."),
				[]byte("##.."),
			},
			true,
		},
	}

	for i, tt := range tests {
		tet, err := ValidateAndCreateTetromino(tt.block, i)
		if tt.wantError && err == nil {
			t.Errorf("Test %d: expected error, got none", i)
		}
		if !tt.wantError && err != nil {
			t.Errorf("Test %d: unexpected error: %v", i, err)
		}
		if !tt.wantError && tet != nil {
			if len(tet.Points) != 4 || tet.Letter != rune('A'+i) {
				t.Errorf("Test %d: invalid tetromino created", i)
			}
		}
	}
}
