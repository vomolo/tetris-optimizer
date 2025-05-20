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
		t.Fatalf("createTestTetromino failed: %v", err)
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
			wantErrMsg: "no tetrominos provided",
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
		t.Fatalf("createTestTetromino failed: %v", err)
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
			wantErrMsg: "not a repetitive case",
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
