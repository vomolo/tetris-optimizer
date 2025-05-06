package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] != "" && filepath.Ext(os.Args[1]) == ".txt" {
		fmt.Println("Valid .txt file provided")
	} else {
		fmt.Println("Usage: go run . <filename>.txt")
	}
}