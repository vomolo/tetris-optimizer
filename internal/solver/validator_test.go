package solver

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_validateAndSolveContent(t *testing.T) {
	// Setup test files
	setupTestFilesForValidateAndSolve(t)
	defer cleanupTestFiles(t)

	tests := []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid single square tetromino",
			filename: "valid_single.txt",
			want:     "AA\nAA",
			wantErr:  false,
		},
		{
			name:     "valid multiple tetrominos",
			filename: "valid_multiple.txt",
			want:     "AAC.\nAAC.\nBBC.\nBBC.",
			wantErr:  false,
		},
		{
			name:     "invalid file - incorrect block size",
			filename: "invalid_block_size.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "invalid file - invalid character",
			filename: "invalid_char.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "invalid file - disconnected tetromino",
			filename: "disconnected.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "invalid file - missing newline between tetrominos",
			filename: "missing_newline.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty file",
			filename: "empty.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "trailing newlines",
			filename: "trailing_newlines.txt",
			want:     "AA\nAA",
			wantErr:  false,
		},
		{
			name:     "many tetrominos",
			filename: "many_tetrominos.txt",
			want:     "AABBCC\nAABBCC\nDDEEFF\nDDEEFF\nGGHHII\nGGHHII",
			wantErr:  false,
		},
		{
			name:     "single with empty lines",
			filename: "single_with_empty.txt",
			want:     "AA\nAA",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndSolveContent(filepath.Join("testfiles", tt.filename))
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndSolveContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateAndSolveContent() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func setupTestFilesForValidateAndSolve(t *testing.T) {
	// Create testfiles directory if it doesn't exist
	if err := os.MkdirAll("testfiles", 0755); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}

	// Define test files with proper formatting
	testFiles := map[string]string{
		"valid_single.txt": "##..\n##..\n....\n....\n",
		"valid_multiple.txt": "##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"#...\n#...\n#...\n#...\n",
		"invalid_block_size.txt": "##.\n##.\n...\n...\n",
		"invalid_char.txt":       "##X.\n##..\n....\n....\n",
		"disconnected.txt":       "#.#.\n....\n#.#.\n....\n",
		"missing_newline.txt": "##..\n##..\n....\n...." +
			"##..\n##..\n....\n....\n",
		"empty.txt":             "",
		"trailing_newlines.txt": "##..\n##..\n....\n....\n\n\n\n",
		"many_tetrominos.txt": "##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n\n" +
			"##..\n##..\n....\n....\n",
		"single_with_empty.txt": "##..\n##..\n....\n....\n\n\n\n\n",
	}

	// Create test files
	for filename, content := range testFiles {
		path := filepath.Join("testfiles", filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal("Failed to create test file:", err)
		}
	}
}

func cleanupTestFiles(t *testing.T) {
	files, err := os.ReadDir("testfiles")
	if err != nil && !os.IsNotExist(err) {
		t.Log("Warning: could not read testfiles directory:", err)
		return
	}

	for _, file := range files {
		if err := os.Remove(filepath.Join("testfiles", file.Name())); err != nil {
			t.Log("Warning: could not remove test file:", err)
		}
	}

	if err := os.Remove("testfiles"); err != nil && !os.IsNotExist(err) {
		t.Log("Warning: could not remove testfiles directory:", err)
	}
}
