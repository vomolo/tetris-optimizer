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

	if err := functions.Validate(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File is valid and optimized for processing")
}