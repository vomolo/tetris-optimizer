package solver

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a < b", 3, 5, 3},
		{"a > b", 5, 3, 3},
		{"a == b", 4, 4, 4},
		{"negative numbers", -2, -1, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.a, tt.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createTetromino(letter rune, points [][2]int, width, height int) *Tetromino {
	t := &Tetromino{
		Letter: letter,
		Width:  width,
		Height: height,
	}
	for _, p := range points {
		t.Points = append(t.Points, Point{X: p[0], Y: p[1]})
	}
	return t
}

func TestSolveTetrominos(t *testing.T) {
	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantErr    bool
	}{
		{
			name:       "empty input",
			tetrominos: []*Tetromino{},
			wantErr:    true,
		},
		{
			name: "single square tetromino",
			tetrominos: []*Tetromino{
				createTetromino('A', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
			},
			wantErr: false,
		},
		{
			name: "two simple tetrominos",
			tetrominos: []*Tetromino{
				createTetromino('A', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
				createTetromino('B', [][2]int{{0, 0}, {0, 1}, {0, 2}, {1, 2}}, 2, 3),
			},
			wantErr: false,
		},
		{
			name: "truly impossible configuration",
			tetrominos: []*Tetromino{
				createTetromino('A', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
				createTetromino('B', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
				createTetromino('C', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
				createTetromino('D', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
				createTetromino('E', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
			},
			wantErr: false,
		},
		{
			name: "complex possible configuration",
			tetrominos: []*Tetromino{
				createTetromino('A', [][2]int{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, 3, 2),
				createTetromino('B', [][2]int{{0, 1}, {1, 1}, {1, 0}, {2, 0}}, 3, 2),
				createTetromino('C', [][2]int{{0, 0}, {0, 1}, {0, 2}, {1, 2}}, 2, 3),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SolveTetrominos(tt.tetrominos)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: SolveTetrominos() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestSolveTetrominosMaxSize(t *testing.T) {
	tetrominos := []*Tetromino{
		createTetromino('A', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
		createTetromino('B', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
		createTetromino('C', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
		createTetromino('D', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
	}
	_, err := SolveTetrominos(tetrominos)
	if err != nil {
		t.Errorf("SolveTetrominos() with max size board failed: %v", err)
	}
}

func TestSolveTetrominosRectangular(t *testing.T) {
	tetrominos := []*Tetromino{
		createTetromino('A', [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, 2, 2),
		createTetromino('B', [][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 4, 1),
	}
	result, err := SolveTetrominos(tetrominos)
	if err != nil {
		t.Errorf("SolveTetrominos() with rectangular board failed: %v", err)
	}
	expected := "AA..\nAA..\nBBBB\n....\n...."
	if result != expected {
		t.Errorf("SolveTetrominos() = \n%v\n, expected \n%v", result, expected)
	}
}
