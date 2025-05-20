package solver

import (
	"os"
	"path/filepath"
	
	"testing"
)

func TestValidate(t *testing.T) {
	// Create a temporary testfiles directory
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "testfiles"), 0755)

	// Create test file
	testFileContent := "##..\n##..\n....\n...."
	testFilePath := filepath.Join(tmpDir, "testfiles", "test.txt")
	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		filename    string
		wantErr     bool
		wantErrMsg  string
		wantBoard   string
	}{
		{
			name:      "ValidFile",
			filename:  "test.txt",
			wantBoard: "AA\nAA\n",
		},
		{
			name:       "InvalidExtension",
			filename:   "test.dat",
			wantErr:    true,
			wantErrMsg: "file must have .txt extension",
		},
		{
			name:       "NonExistentFile",
			filename:   "nonexistent.txt",
			wantErr:    true,
			wantErrMsg: "file does not exist in directory",
		},
		{
			name:       "DirectoryTraversal",
			filename:   "../test.txt",
			wantErr:    true,
			wantErrMsg: "invalid file path: attempted directory traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to temp directory for test
			origDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(origDir)

			got, err := Validate(tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil; want %q", tt.wantErrMsg)
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("Validate() error = %q; want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Validate() error = %v; want nil", err)
			}
			if !compareBoards(got, tt.wantBoard) {
				t.Errorf("Validate() = %q; want %q", got, tt.wantBoard)
			}
		})
	}
}
