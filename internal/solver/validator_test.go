package solver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	// Setup test files
	createTestFiles(t)
	defer cleanupTestFiles(t)

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid single square tetromino",
			args:    args{filename: "valid_single.txt"},
			want:    "AA\nAA",
			wantErr: false,
		},
		{
			name:    "valid multiple tetrominos",
			args:    args{filename: "valid_multiple.txt"},
			want:    "AABB\nAABB\nBB..\nBB..",
			wantErr: false,
		},
		{
			name:    "invalid file - wrong extension",
			args:    args{filename: "wrong_extension.dat"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - doesn't exist",
			args:    args{filename: "nonexistent.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - empty file",
			args:    args{filename: "empty.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - incorrect block size",
			args:    args{filename: "invalid_block_size.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - invalid character",
			args:    args{filename: "invalid_char.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - disconnected tetromino",
			args:    args{filename: "disconnected.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid file - missing newline between tetrominos",
			args:    args{filename: "missing_newline.txt"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "valid file - maximum size board",
			args:    args{filename: "max_size.txt"},
			want:    "AAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions to create and clean up test files
func createTestFiles(t *testing.T) {
	testFiles := map[string]string{
		"valid_single.txt": `##..
##..
....
....`,

		"valid_multiple.txt": `##..
##..
....
....

##..
##..
....
....

#...
#...
#...
#...`,

		"empty.txt": "",

		"invalid_block_size.txt": `##.
##.
...
...`,

		"invalid_char.txt": `##X.
##..
....
....`,

		"disconnected.txt": `#.#.
....
#.#.
....`,

		"missing_newline.txt": `##..
##..
....
....##..
##..
....
....`,

		"max_size.txt": `##########
##########
##########
##########
##########
##########
##########
##########
##########
##########`,
	}

	for filename, content := range testFiles {
		path := filepath.Join("testfiles", filename)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}
}

func cleanupTestFiles(t *testing.T) {
	files := []string{
		"valid_single.txt",
		"valid_multiple.txt",
		"empty.txt",
		"invalid_block_size.txt",
		"invalid_char.txt",
		"disconnected.txt",
		"missing_newline.txt",
		"max_size.txt",
	}

	for _, filename := range files {
		path := filepath.Join("testfiles", filename)
		os.Remove(path)
	}
}
