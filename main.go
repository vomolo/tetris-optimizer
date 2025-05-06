package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// Check if the correct number of arguments is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Check .txt extension
	if filepath.Ext(filename) != ".txt" {
		fmt.Println("Error: File must have .txt extension")
		os.Exit(1)
	}

	// Automatically prepend "tetris_files/" to the filename
	fullPath := filepath.Join("tetris_files", filename)

	// Check if file exists in tetris_files
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found in tetris_files directory\n", filename)
		os.Exit(1)
	}

	fmt.Println("File is valid and ready for processing!")
	fmt.Println("Full path:", fullPath)
}