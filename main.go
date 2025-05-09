package main

import (
	"fmt"
	"os"
	"tetris_optimizer/internal/solver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <filename>")
		os.Exit(1)
	}

	solution, err := solver.Validate(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR")
		os.Exit(1)
	}

	fmt.Println(solution)
}
