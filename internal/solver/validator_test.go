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
	defer cleanupTestFilesForValidate(t)

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
func createTestFiles(_ *testing.T) error {
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

func cleanupTestFilesForValidate(_ *testing.T) {
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

func Test_validateStructure(t *testing.T) {
	// Setup test environment
	setupTestFiles(t)
	defer cleanupTestFiles(t)

	tests := []struct {
		name     string
		fullPath string
		wantErr  bool
	}{
		{
			name:     "valid txt file",
			fullPath: filepath.Join("testfiles", "valid.txt"),
			wantErr:  false,
		},
		{
			name:     "wrong file extension",
			fullPath: filepath.Join("testfiles", "invalid.dat"),
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			fullPath: filepath.Join("testfiles", "nonexistent.txt"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStructure(tt.fullPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupTestFiles(t *testing.T) {
	// Create test directory
	if err := os.MkdirAll("testfiles", 0755); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}

	// Create test files
	testFiles := []struct {
		name    string
		content string
	}{
		{"valid.txt", "test content"},
		{"invalid.dat", "test content"},
	}

	for _, tf := range testFiles {
		path := filepath.Join("testfiles", tf.name)
		if err := os.WriteFile(path, []byte(tf.content), 0644); err != nil {
			t.Fatal("Failed to create test file:", err)
		}
	}
}

func cleanupTestFiles(t *testing.T) {
	// Remove test files
	os.Remove(filepath.Join("testfiles", "valid.txt"))
	os.Remove(filepath.Join("testfiles", "invalid.dat"))

	// Remove test directory
	os.Remove("testfiles")
}
