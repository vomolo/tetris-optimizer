package solver

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small board", 4, 4},
		{"rectangle board", 5, 3},
		{"single row", 10, 1},
		{"single column", 1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard(tt.width, tt.height)
			if b.Width != tt.width || b.Height != tt.height {
				t.Errorf("Expected board dimensions %dx%d, got %dx%d", tt.width, tt.height, b.Width, b.Height)
			}
			if len(b.Grid) != tt.height {
				t.Errorf("Expected grid height %d, got %d", tt.height, len(b.Grid))
			}
			for _, row := range b.Grid {
				if len(row) != tt.width {
					t.Errorf("Expected row width %d, got %d", tt.width, len(row))
				}
				for _, cell := range row {
					if cell != 0 {
						t.Error("New board should be initialized with zeros")
					}
				}
			}
			if b.Placed != 0 {
				t.Error("New board should have zero placed pieces")
			}
		})
	}
}

func TestCanPlace(t *testing.T) {
	b := NewBoard(4, 4)
	tetromino := &Tetromino{
		Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, // 2x2 square
		Letter: 'A',
	}

	tests := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{"top-left", 0, 0, true},
		{"middle", 1, 1, true},
		{"out of bounds right", 3, 0, false},
		{"out of bounds bottom", 0, 3, false},
		{"negative x", -1, 0, false},
		{"negative y", 0, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := b.canPlace(tetromino, tt.x, tt.y); result != tt.expected {
				t.Errorf("canPlace(%d, %d) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}

	// Test with occupied cells
	b.place(tetromino, 0, 0)
	tests = []struct {
		name     string
		x, y     int
		expected bool
	}{
		{"overlap placed", 0, 0, false},
		{"adjacent right", 2, 0, true},
		{"adjacent down", 0, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := b.canPlace(tetromino, tt.x, tt.y); result != tt.expected {
				t.Errorf("canPlace(%d, %d) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

