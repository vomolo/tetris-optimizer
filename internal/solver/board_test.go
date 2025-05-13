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
