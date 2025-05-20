package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
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
		name       string
		args       []string
		wantOutput string
		wantExit   int
	}{
		{
			name:       "ValidFile",
			args:       []string{"program", "test.txt"},
			wantOutput: "AA\nAA\n",
			wantExit:   0,
		},
		{
			name:       "NoArguments",
			args:       []string{"program"},
			wantOutput: "Usage: go run main.go <filename>\n",
			wantExit:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to temp directory
			origDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(origDir)

			// Redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = oldStdout }()

			// Redirect stderr
			oldStderr := os.Stderr
			rErr, wErr, _ := os.Pipe()
			os.Stderr = wErr
			defer func() { os.Stderr = oldStderr }()

			// Set args and run main
			os.Args = tt.args
			var exitCode int
			func() {
				defer func() {
					if r := recover(); r != nil {
						if code, ok := r.(int); ok {
							exitCode = code
						}
					}
				}()
				main()
			}()

			// Capture output
			w.Close()
			wErr.Close()
			var stdoutBuf, stderrBuf bytes.Buffer
			_, _ = stdoutBuf.ReadFrom(r)
			_, _ = stderrBuf.ReadFrom(rErr)

			output := stdoutBuf.String() + stderrBuf.String()
			if output != tt.wantOutput {
				t.Errorf("main() output = %q; want %q", output, tt.wantOutput)
			}
			if exitCode != tt.wantExit {
				t.Errorf("main() exit code = %d; want %d", exitCode, tt.wantExit)
			}
		})
	}
}