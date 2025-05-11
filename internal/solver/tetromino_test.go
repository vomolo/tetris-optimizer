package solver

import (
	"reflect"
	"testing"
)

func Test_validateAndCreateTetromino(t *testing.T) {
	type args struct {
		block       [][]byte
		blockNumber int
	}
	tests := []struct {
		name    string
		args    args
		want    *Tetromino
		wantErr bool
	}{
		{
			name: "valid square tetromino",
			args: args{
				block: [][]byte{
					{'#', '#', '.', '.'},
					{'#', '#', '.', '.'},
					{'.', '.', '.', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 0,
			},
			want: &Tetromino{
				Points: []Point{
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
				},
				Letter: 'A',
				Width:  2,
				Height: 2,
			},
			wantErr: false,
		},
		{
			name: "valid line tetromino",
			args: args{
				block: [][]byte{
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
				},
				blockNumber: 1,
			},
			want: &Tetromino{
				Points: []Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 0, Y: 2},
					{X: 0, Y: 3},
				},
				Letter: 'B',
				Width:  1,
				Height: 4,
			},
			wantErr: false,
		},
		{
			name: "valid L-shaped tetromino",
			args: args{
				block: [][]byte{
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'#', '#', '.', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 2,
			},
			want: &Tetromino{
				Points: []Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 0, Y: 2},
					{X: 1, Y: 2},
				},
				Letter: 'C',
				Width:  2,
				Height: 3,
			},
			wantErr: false,
		},
		{
			name: "invalid tetromino - too few #",
			args: args{
				block: [][]byte{
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 3,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid tetromino - too many #",
			args: args{
				block: [][]byte{
					{'#', '#', '.', '.'},
					{'#', '#', '.', '.'},
					{'#', '#', '.', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 4,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid tetromino - disconnected",
			args: args{
				block: [][]byte{
					{'#', '.', '#', '.'},
					{'.', '.', '.', '.'},
					{'#', '.', '#', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 5,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid character in block",
			args: args{
				block: [][]byte{
					{'#', '.', 'X', '.'},
					{'#', '.', '.', '.'},
					{'#', '#', '.', '.'},
					{'.', '.', '.', '.'},
				},
				blockNumber: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid block size - too few rows",
			args: args{
				block: [][]byte{
					{'#', '.', '.', '.'},
					{'#', '.', '.', '.'},
					{'#', '#', '.', '.'},
				},
				blockNumber: 7,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid block size - too few columns",
			args: args{
				block: [][]byte{
					{'#', '.', '.'},
					{'#', '.', '.'},
					{'#', '#', '.'},
					{'.', '.', '.'},
				},
				blockNumber: 8,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndCreateTetromino(tt.args.block, tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndCreateTetromino() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateAndCreateTetromino() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidTetromino(t *testing.T) {
	type args struct {
		points [4]Point
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid square tetromino",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
				},
			},
			want: true,
		},
		{
			name: "valid line tetromino (horizontal)",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 2, Y: 0},
					{X: 3, Y: 0},
				},
			},
			want: true,
		},
		{
			name: "valid line tetromino (vertical)",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 0, Y: 2},
					{X: 0, Y: 3},
				},
			},
			want: true,
		},
		{
			name: "valid L-shaped tetromino",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 0, Y: 2},
					{X: 1, Y: 2},
				},
			},
			want: true,
		},
		{
			name: "valid T-shaped tetromino",
			args: args{
				points: [4]Point{
					{X: 1, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
					{X: 2, Y: 1},
				},
			},
			want: true,
		},
		{
			name: "valid S-shaped tetromino",
			args: args{
				points: [4]Point{
					{X: 1, Y: 0},
					{X: 2, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
				},
			},
			want: true,
		},
		{
			name: "invalid tetromino - disconnected points",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 2, Y: 0},
					{X: 0, Y: 2},
					{X: 2, Y: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid tetromino - three connected with one separate",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 2, Y: 0},
					{X: 0, Y: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid tetromino - diagonal connection only",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
					{X: 2, Y: 2},
					{X: 3, Y: 3},
				},
			},
			want: false,
		},
		{
			name: "invalid tetromino - duplicate points",
			args: args{
				points: [4]Point{
					{X: 0, Y: 0},
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 2, Y: 0},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTetromino(tt.args.points); got != tt.want {
				t.Errorf("isValidTetromino() = %v, want %v", got, tt.want)
			}
		})
	}
}
