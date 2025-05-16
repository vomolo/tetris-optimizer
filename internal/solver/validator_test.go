package solver

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_validateAndSolve(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			name:    "valid single square tetromino",
			content: "##..\n##..\n....\n....\n",
			want:    "AA\nAA",
			wantErr: false,
		},
		{
			name: "valid multiple tetrominos",
			content: "##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"#...\n#...\n#...\n#...\n",
			want:    "AAC\nAAC\nBBC\nBBC",
			wantErr: false,
		},
		{
			name:    "invalid block size",
			content: "##.\n##.\n...\n...\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid character",
			content: "##X.\n##..\n....\n....\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "disconnected tetromino",
			content: "#.#.\n....\n#.#.\n....\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing newline between tetrominos",
			content: "##..\n##..\n....\n...." +
				"##..\n##..\n....\n....\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty content",
			content: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "trailing newlines",
			content: "##..\n##..\n....\n....\n\n\n\n",
			want:    "",
			wantErr: true,
		},
		{
			name: "many tetrominos",
			content: "##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n\n" +
				"##..\n##..\n....\n....\n",
			want:    "AABBCCDDEE\nAABBCCDDEE\nFFGGHHIIJJ\nFFGGHHIIJJ",
			wantErr: false,
		},
		{
			name:    "single with empty lines",
			content: "##..\n##..\n....\n....\n\n\n\n\n",
			want:    "AA\nAA",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndSolve(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndSolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateAndSolve() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func Test_validateAndSolveContent_singleCase(t *testing.T) {
	setupTestFilesForValidateAndSolve(t)
	defer cleanupTestFiles(t)

	want := "AA\nAA"
	got, err := validateAndSolveContent(filepath.Join("testfiles", "valid_single.txt"))
	if err != nil {
		t.Errorf("validateAndSolveContent() unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("validateAndSolveContent() =\n%v\nwant\n%v", got, want)
	}
}

func setupTestFilesForValidateAndSolve(t *testing.T) {
	if err := os.MkdirAll("testfiles", 0755); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}
	content := "##..\n##..\n....\n....\n"
	path := filepath.Join("testfiles", "valid_single.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal("Failed to create test file:", err)
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
