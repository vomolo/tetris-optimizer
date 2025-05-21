package solver

import (
	"reflect"
	"testing"
)

func TestValidateAndCreateTetromino(t *testing.T) {
	tests := []struct {
		name        string
		block       [][]byte
		wantErr     bool
		wantErrMsg  string
		wantPoints  []Point
		wantWidth   int
		wantHeight  int
	}{
		{
			name: "ValidSquare",
			block: [][]byte{
				[]byte("##.."),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			wantPoints: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			wantWidth:  2,
			wantHeight: 2,
		},
		{
			name: "InvalidRowCount",
			block: [][]byte{
				[]byte("##.."),
				[]byte("##.."),
			},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name: "InvalidColumnCount",
			block: [][]byte{
				[]byte("###"),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name: "TooManyBlocks",
			block: [][]byte{
				[]byte("###."),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name: "InvalidCharacter",
			block: [][]byte{
				[]byte("#X.."),
				[]byte("##.."),
				[]byte("...."),
				[]byte("...."),
			},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
		{
			name: "DisconnectedTetromino",
			block: [][]byte{
				[]byte("#.##"),
				[]byte("...."),
				[]byte("...."),
				[]byte("...."),
			},
			wantErr:    true,
			wantErrMsg: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndCreateTetromino(tt.block, 0)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAndCreateTetromino() error = nil; want %q", tt.wantErrMsg)
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("ValidateAndCreateTetromino() error = %q; want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("ValidateAndCreateTetromino() error = %v; want nil", err)
			}
			if !reflect.DeepEqual(got.Points, tt.wantPoints) {
				t.Errorf("ValidateAndCreateTetromino() Points = %v; want %v", got.Points, tt.wantPoints)
			}
			if got.Width != tt.wantWidth {
				t.Errorf("ValidateAndCreateTetromino() Width = %d; want %d", got.Width, tt.wantWidth)
			}
			if got.Height != tt.wantHeight {
				t.Errorf("ValidateAndCreateTetromino() Height = %d; want %d", got.Height, tt.wantHeight)
			}
		})
	}
}

func TestIsValidTetromino(t *testing.T) {
	tests := []struct {
		name   string
		points [4]Point
		want   bool
	}{
		{
			name:   "ValidSquare",
			points: [4]Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			want:   true,
		},
		{
			name:   "Disconnected",
			points: [4]Point{{0, 0}, {1, 0}, {0, 2}, {1, 2}},
			want:   false,
		},
		{
			name:   "DuplicatePoints",
			points: [4]Point{{0, 0}, {0, 0}, {0, 1}, {1, 1}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTetromino(tt.points); got != tt.want {
				t.Errorf("isValidTetromino() = %v; want %v", got, tt.want)
			}
		})
	}
}