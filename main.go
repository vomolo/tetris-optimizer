package main

import (
	"fmt"
	"os"
	"tet/functions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <filename>")
		fmt.Println("Note: The file should be in the tetris_files directory")
		os.Exit(1)
	}

	filename := os.Args[1]

	fullPath, err := functions.RunInitialChecks(filename)
	if err != nil {
		fmt.Println("Validation Error:", err)
		os.Exit(1)
	}

	fmt.Println("All initial checks passed successfully!")
	fmt.Printf("Processing file: %s\n", fullPath)
	
}