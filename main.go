package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := "input.txt"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ERROR opening file")
		return
	}
	defer file.Close()
}
