package solver

import (
	"reflect"
	"testing"
)

func TestTetrominoValidation(t *testing.T) {
	t.Run("validateAndCreateTetromino", testValidateAndCreateTetromino)
	t.Run("isValidTetromino", testIsValidTetromino)
}

func testValidateAndCreateTetromino(t *testing.T) {
	tests := []struct {
		name      string
		block     [][]byte
		blockNum  int
		want      *Tetromino
		wantError bool
	}{
		{
			name: "valid square",
			block: [][]byte{
				{'#', '#', '.', '.'},
				{'#', '#', '.', '.'},
				{'.', '.', '.', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum: 0,
			want: &Tetromino{
				Points: []Point{
					{0, 0}, {1, 0},
					{0, 1}, {1, 1},
				},
				Letter: 'A',
				Width:  2,
				Height: 2,
			},
		},
		{
			name: "valid vertical line",
			block: [][]byte{
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
			},
			blockNum: 1,
			want: &Tetromino{
				Points: []Point{
					{0, 0}, {0, 1},
					{0, 2}, {0, 3},
				},
				Letter: 'B',
				Width:  1,
				Height: 4,
			},
		},
		{
			name: "valid L-shape",
			block: [][]byte{
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'#', '#', '.', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum: 2,
			want: &Tetromino{
				Points: []Point{
					{0, 0}, {0, 1},
					{0, 2}, {1, 2},
				},
				Letter: 'C',
				Width:  2,
				Height: 3,
			},
		},
		{
			name: "too few # characters",
			block: [][]byte{
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum:  3,
			wantError: true,
		},
		{
			name: "too many # characters",
			block: [][]byte{
				{'#', '#', '.', '.'},
				{'#', '#', '.', '.'},
				{'#', '#', '.', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum:  4,
			wantError: true,
		},
		{
			name: "disconnected points",
			block: [][]byte{
				{'#', '.', '#', '.'},
				{'.', '.', '.', '.'},
				{'#', '.', '#', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum:  5,
			wantError: true,
		},
		{
			name: "invalid character",
			block: [][]byte{
				{'#', '.', 'X', '.'},
				{'#', '.', '.', '.'},
				{'#', '#', '.', '.'},
				{'.', '.', '.', '.'},
			},
			blockNum:  6,
			wantError: true,
		},
		{
			name: "too few rows",
			block: [][]byte{
				{'#', '.', '.', '.'},
				{'#', '.', '.', '.'},
				{'#', '#', '.', '.'},
			},
			blockNum:  7,
			wantError: true,
		},
		{
			name: "too few columns",
			block: [][]byte{
				{'#', '.', '.'},
				{'#', '.', '.'},
				{'#', '#', '.'},
				{'.', '.', '.'},
			},
			blockNum:  8,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndCreateTetromino(tt.block, tt.blockNum)
			if (err != nil) != tt.wantError {
				t.Errorf("got error %v, want error %v", err != nil, tt.wantError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func testIsValidTetromino(t *testing.T) {
	tests := []struct {
		name   string
		points [4]Point
		want   bool
	}{
		{
			name: "valid square",
			points: [4]Point{
				{0, 0}, {1, 0},
				{0, 1}, {1, 1},
			},
			want: true,
		},
		{
			name: "valid horizontal line",
			points: [4]Point{
				{0, 0}, {1, 0},
				{2, 0}, {3, 0},
			},
			want: true,
		},
		{
			name: "valid vertical line",
			points: [4]Point{
				{0, 0}, {0, 1},
				{0, 2}, {0, 3},
			},
			want: true,
		},
		{
			name: "valid L-shape",
			points: [4]Point{
				{0, 0}, {0, 1},
				{0, 2}, {1, 2},
			},
			want: true,
		},
		{
			name: "valid T-shape",
			points: [4]Point{
				{1, 0}, {0, 1},
				{1, 1}, {2, 1},
			},
			want: true,
		},
		{
			name: "valid S-shape",
			points: [4]Point{
				{1, 0}, {2, 0},
				{0, 1}, {1, 1},
			},
			want: true,
		},
		{
			name: "disconnected points",
			points: [4]Point{
				{0, 0}, {2, 0},
				{0, 2}, {2, 2},
			},
			want: false,
		},
		{
			name: "three connected + one separate",
			points: [4]Point{
				{0, 0}, {1, 0},
				{2, 0}, {0, 2},
			},
			want: false,
		},
		{
			name: "diagonal connections only",
			points: [4]Point{
				{0, 0}, {1, 1},
				{2, 2}, {3, 3},
			},
			want: false,
		},
		{
			name: "duplicate points",
			points: [4]Point{
				{0, 0}, {0, 0},
				{1, 0}, {2, 0},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidTetromino(tt.points)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
