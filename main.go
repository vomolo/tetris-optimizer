package main

import (
	"bufio"
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

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if lineCount%5 != 0 && line == "" {
			fmt.Println("ERROR: Unwanted empty line:", line)
			return
		}

	}
	fmt.Println("Success")
}
