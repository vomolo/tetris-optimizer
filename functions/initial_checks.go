package functions

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckTxtExtension verifies the file has .txt extension
func CheckTxtExtension(filename string) error {
	if filepath.Ext(filename) != ".txt" {
		return fmt.Errorf("file must have .txt extension")
	}
	return nil
}

// CheckInTetrisFiles verifies the file is in tetris_files directory
func CheckInTetrisFiles(fullPath string) error {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	requiredDir := "tetris_files"
	dir := filepath.Base(filepath.Dir(absPath))
	if dir != requiredDir {
		return fmt.Errorf("file must be in '%s' directory", requiredDir)
	}
	return nil
}

// CheckFileExists verifies the file exists
func CheckFileExists(fullPath string) error {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist", fullPath)
	}
	return nil
}

// RunInitialChecks combines all validation checks
func RunInitialChecks(filename string) (string, error) {
	// Automatically prepend "tetris_files/" to the filename
	fullPath := filepath.Join("tetris_files", filename)

	checks := []struct {
		name string
		fn   func(string) error
	}{
		{"extension check", CheckTxtExtension},
		{"directory check", func(string) error { return CheckInTetrisFiles(fullPath) }},
		{"existence check", func(string) error { return CheckFileExists(fullPath) }},
	}

	for _, check := range checks {
		if err := check.fn(filename); err != nil {
			return "", fmt.Errorf("%s failed: %v", check.name, err)
		}
	}

	return fullPath, nil
}