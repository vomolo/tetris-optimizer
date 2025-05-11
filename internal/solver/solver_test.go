package solver

import "testing"

func TestSolveTetrominos(t *testing.T) {
	type args struct {
		tetrominos []*Tetromino
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "single square tetromino",
			args: args{
				tetrominos: []*Tetromino{
					{
						Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
						Letter: 'A',
						Width:  2,
						Height: 2,
					},
				},
			},
			want:    "AA\nAA",
			wantErr: false,
		},
		{
			name: "two square tetrominos",
			args: args{
				tetrominos: []*Tetromino{
					{
						Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
						Letter: 'A',
						Width:  2,
						Height: 2,
					},
					{
						Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
						Letter: 'B',
						Width:  2,
						Height: 2,
					},
				},
			},
			want:    "AABB\nAABB\n....\n....",
			wantErr: false,
		},
		{
			name: "square and line tetrominos",
			args: args{
				tetrominos: []*Tetromino{
					{
						Points: []Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
						Letter: 'A',
						Width:  2,
						Height: 2,
					},
					{
						Points: []Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
						Letter: 'B',
						Width:  1,
						Height: 4,
					},
				},
			},
			want:    "AAB.\nAAB.\n..B.\n..B.",
			wantErr: false,
		},
		{
			name: "three L-shaped tetrominos",
			args: args{
				tetrominos: []*Tetromino{
					{
						Points: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
						Letter: 'A',
						Width:  2,
						Height: 3,
					},
					{
						Points: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
						Letter: 'B',
						Width:  2,
						Height: 3,
					},
					{
						Points: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
						Letter: 'C',
						Width:  2,
						Height: 3,
					},
				},
			},
			want:    "A..B.\nA.CB.\nAACBB\n..CC.\n.....",
			wantErr: false,
		},
		{
			name: "empty tetrominos list",
			args: args{
				tetrominos: []*Tetromino{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "impossible configuration",
			args: args{
				tetrominos: []*Tetromino{
					{
						Points: []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
						Letter: 'A',
						Width:  4,
						Height: 1,
					},
					{
						Points: []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
						Letter: 'B',
						Width:  4,
						Height: 1,
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SolveTetrominos(tt.args.tetrominos)
			if (err != nil) != tt.wantErr {
				t.Errorf("SolveTetrominos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SolveTetrominos() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
