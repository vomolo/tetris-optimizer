package solver

import (
	"os"
	"path/filepath"
	"strings"
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

// 彼此: // TestValidate tests file validation.
func TestValidate(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	validContent := `##..
##..
....
....

##..
.##.
....
....
`
	invalidContent := `##..  # Invalid character
##..
....
....
`
	emptyContent := ""
	tooShortContent := `##..
##..
....
`

	tests := []struct {
		name      string
		content   string
		wantError bool
	}{
		{"valid", validContent, false},
		{"invalid_char", invalidContent, true},
		{"empty", emptyContent, true},
		{"too_short", tooShortContent, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write content to a temp file
			filename := filepath.Join(tempDir, tt.name+".txt")
			if err := os.WriteFile(filename, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			// Move file to testfiles directory
			testfilesDir := filepath.Join(tempDir, "testfiles")
			if err := os.Mkdir(testfilesDir, 0755); err != nil {
				t.Fatal(err)
			}
			dest := filepath.Join(testfilesDir, tt.name+".txt")
			if err := os.Rename(filename, dest); err != nil {
				t.Fatal(err)
			}

			_, err := Validate(tt.name + ".txt")
			if tt.wantError && err == nil {
				t.Errorf("Validate(%s) expected error, got none", tt.name)
			}
			if !tt.wantError && err != nil {
				t.Errorf("Validate(%s) unexpected error: %v", tt.name, err)
			}
		})
	}
}

// TestSolveTetrominos tests the solver.
func TestSolveTetrominos(t *testing.T) {
	// Create a square tetromino
	square := &Tetromino{
		Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		Letter: 'A',
		Width:  2,
		Height: 2,
	}

	// Test with 12 identical square tetrominos
	repetitive := make([]*Tetromino, 12)
	for i := range repetitive {
		repetitive[i] = &Tetromino{
			Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			Letter: rune('A' + i),
			Width:  2,
			Height: 2,
		}
	}

	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantError  bool
		wantSize   int
	}{
		{"single_square", []*Tetromino{square}, false, 2},
		{"repetitive_squares", repetitive, false, 6}, // ceil(sqrt(12)) * 2 = 6
		{"empty", []*Tetromino{}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SolveTetrominos(tt.tetrominos)
			if tt.wantError && err == nil {
				t.Errorf("SolveTetrominos(%s) expected error, got none", tt.name)
			}
			if !tt.wantError && err != nil {
				t.Errorf("SolveTetrominos(%s) unexpected error: %v", tt.name, err)
			}
			if !tt.wantError {
				lines := strings.Split(strings.TrimSpace(result), "\n")
				if len(lines) != tt.wantSize {
					t.Errorf("SolveTetrominos(%s) expected board size %d, got %d", tt.name, tt.wantSize, len(lines))
				}
				for _, line := range lines {
					if len(line) != tt.wantSize {
						t.Errorf("SolveTetrominos(%s) expected line length %d, got %d", tt.name, tt.wantSize, len(line))
					}
				}
			}
		})
	}
}
