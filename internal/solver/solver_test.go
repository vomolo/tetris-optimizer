package solver

import (
	"testing"
)

func TestSolveTetrominos(t *testing.T) {
	tetromino, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 0)
	if err != nil {
		t.Fatalf("ERROR")
	}

	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantErr    bool
		wantErrMsg string
		wantBoard  string
	}{
		{
			name:       "EmptyInput",
			tetrominos: []*Tetromino{},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name:       "SingleSquareTetromino",
			tetrominos: []*Tetromino{tetromino},
			wantBoard:  "AA\nAA\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SolveTetrominos(tt.tetrominos)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SolveTetrominos() error = nil; want %q", tt.wantErrMsg)
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("SolveTetrominos() error = %q; want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("SolveTetrominos() error = %v; want nil", err)
			}
			if !compareBoards(got, tt.wantBoard) {
				t.Errorf("SolveTetrominos() = %q; want %q", got, tt.wantBoard)
			}
		})
	}
}

func TestTryOptimizedSquareRepetitiveSolution(t *testing.T) {
	tetromino, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 0)
	if err != nil {
		t.Fatalf("ERROR")
	}

	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantErr    bool
		wantErrMsg string
		wantBoard  string
	}{
		{
			name:       "SingleTetromino",
			tetrominos: []*Tetromino{tetromino},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name: "FiveIdenticalTetrominos",
			tetrominos: []*Tetromino{
				tetromino, tetromino, tetromino, tetromino, tetromino,
			},
			wantBoard: "AAAAAA\nAAAAAA\nAAAA..\nAAAA..\n......\n......",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tryOptimizedSquareRepetitiveSolution(tt.tetrominos)
			if tt.wantErr {
				if err == nil {
					t.Errorf("tryOptimizedSquareRepetitiveSolution() error = nil; want %q", tt.wantErrMsg)
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("tryOptimizedSquareRepetitiveSolution() error = %q; want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("tryOptimizedSquareRepetitiveSolution() error = %v; want nil", err)
			}
			if !compareBoards(got, tt.wantBoard) {
				t.Errorf("tryOptimizedSquareRepetitiveSolution() = %q; want %q", got, tt.wantBoard)
			}
		})
	}
}

func TestGeneralSquareSolver(t *testing.T) {
	tetromino1, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 0)
	if err != nil {
		t.Fatalf("ERROR")
	}
	tetromino2, err := createTestTetromino([]string{
		"#...",
		"###.",
		"....",
		"....",
	}, 1)
	if err != nil {
		t.Fatalf("ERROR")
	}

	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantBoard  string
		wantErr    bool
	}{
		{
			name:       "TwoDifferentTetrominos",
			tetrominos: []*Tetromino{tetromino1, tetromino2},
			wantBoard:  ".AA\nBAA\nBBB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generalSquareSolver(tt.tetrominos)
			if tt.wantErr {
				if err == nil {
					t.Errorf("generalSquareSolver() error = nil; want error")
				}
				return
			}
			if err != nil {
				t.Errorf("generalSquareSolver() error = %v; want nil", err)
			}
			if !compareBoards(got, tt.wantBoard) {
				t.Errorf("generalSquareSolver() = %q; want %q", got, tt.wantBoard)
			}
		})
	}
}

func TestGroupRepetitiveTetrominos(t *testing.T) {
	tetromino1, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 0)
	if err != nil {
		t.Fatalf("ERROR")
	}
	tetromino2, err := createTestTetromino([]string{
		"#...",
		"###.",
		"....",
		"....",
	}, 1)
	if err != nil {
		t.Fatalf("ERROR")
	}

	tests := []struct {
		name       string
		tetrominos []*Tetromino
		wantGroups int
	}{
		{
			name:       "TwoIdentical",
			tetrominos: []*Tetromino{tetromino1, tetromino1},
			wantGroups: 1,
		},
		{
			name:       "TwoDifferent",
			tetrominos: []*Tetromino{tetromino1, tetromino2},
			wantGroups: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := groupRepetitiveTetrominos(tt.tetrominos)
			if len(groups) != tt.wantGroups {
				t.Errorf("groupRepetitiveTetrominos() len = %d; want %d", len(groups), tt.wantGroups)
			}
		})
	}
}

func TestAreTetrominosEqual(t *testing.T) {
	tetromino1, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 0)
	if err != nil {
		t.Fatalf("ERROR")
	}
	tetromino2, err := createTestTetromino([]string{
		"##..",
		"##..",
		"....",
		"....",
	}, 1)
	if err != nil {
		t.Fatalf("ERROR")
	}
	tetromino3, err := createTestTetromino([]string{
		"#...",
		"###.",
		"....",
		"....",
	}, 2)
	if err != nil {
		t.Fatalf("ERROR")
	}

	tests := []struct {
		name   string
		t1, t2 *Tetromino
		want   bool
	}{
		{"IdenticalTetrominos", tetromino1, tetromino2, true},
		{"DifferentTetrominos", tetromino1, tetromino3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areTetrominosEqual(tt.t1, tt.t2); got != tt.want {
				t.Errorf("areTetrominosEqual() = %v; want %v", got, tt.want)
			}
		})
	}
}
