package main

import (
	"fmt"
	"os"
	"tetris_optimizer/functions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: optimizer <filename>")
		os.Exit(1)
	}

	// Get both the formatted output and validation error
	output, err := functions.Validate(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	// Print the formatted output with letter replacements
	fmt.Println("Valid tetromino configuration:")
	fmt.Print(output)
}
