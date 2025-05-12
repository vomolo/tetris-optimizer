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

func TestPlaceAndRemove(t *testing.T) {
	b := NewBoard(4, 4)
	tetromino := &Tetromino{
		Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		Letter: 'B',
	}

	// Test place
	b.place(tetromino, 0, 0)
	if b.Placed != 1 {
		t.Errorf("Expected Placed=1 after placing, got %d", b.Placed)
	}
	for _, p := range tetromino.Points {
		if b.Grid[p.Y][p.X] != tetromino.Letter {
			t.Errorf("Expected cell (%d,%d) to be '%c', got '%c'", p.X, p.Y, tetromino.Letter, b.Grid[p.Y][p.X])
		}
	}

	// Test remove
	b.remove(tetromino, 0, 0)
	if b.Placed != 0 {
		t.Errorf("Expected Placed=0 after removing, got %d", b.Placed)
	}
	for _, p := range tetromino.Points {
		if b.Grid[p.Y][p.X] != 0 {
			t.Errorf("Expected cell (%d,%d) to be empty after remove, got '%c'", p.X, p.Y, b.Grid[p.Y][p.X])
		}
	}
}

func TestBoardToString(t *testing.T) {
	tests := []struct {
		name     string
		board    *Board
		expected string
	}{
		{
			"empty board",
			NewBoard(3, 2),
			"...\n...",
		},
		{
			"partially filled",
			func() *Board {
				b := NewBoard(4, 4)
				t := &Tetromino{
					Points: []Point{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, // T-shape
					Letter: 'T',
				}
				b.place(t, 1, 1)
				return b
			}(),
			"....\n.TTT\n..T.\n....",
		},
		{
			"full row",
			func() *Board {
				b := NewBoard(3, 3)
				t := &Tetromino{
					Points: []Point{{0, 0}, {1, 0}, {2, 0}},
					Letter: 'L',
				}
				b.place(t, 0, 1)
				return b
			}(),
			"...\nLLL\n...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := boardToString(tt.board)
			if result != tt.expected {
				t.Errorf("boardToString() = \n%v\n, expected \n%v", result, tt.expected)
			}
		})
	}
}

