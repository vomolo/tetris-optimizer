package solver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	// Setup test files
	err := createTestFiles(t)
	if err != nil {
		t.Fatalf("Failed to create test files: %v", err)
	}
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
		// ... (keep other test cases the same) ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(filepath.Join("testfiles", tt.args.filename))
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
func createTestFiles(t *testing.T) error {
	// Create testfiles directory if it doesn't exist
	if err := os.MkdirAll("testfiles", 0755); err != nil {
		return err
	}

	testFiles := map[string]string{
		"valid_single.txt": "##..\n##..\n....\n....",
		// ... (keep other test file contents the same) ...
	}

	for filename, content := range testFiles {
		path := filepath.Join("testfiles", filename)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupTestFiles(t *testing.T) {
	files := []string{
		"valid_single.txt",
		// ... (keep other filenames the same) ...
	}

	for _, filename := range files {
		path := filepath.Join("testfiles", filename)
		os.Remove(path)
	}
	// Remove the testfiles directory
	os.Remove("testfiles")
}
