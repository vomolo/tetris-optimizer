package main

import (
	"bufio"
	"fmt"
	"os"
)

// isValidTetromino checks if the given grid represents a valid tetromino.
func isValidTetromino(grid [4][4]rune) bool {
	// Count the number of '#' characters and their positions.
	var positions []struct{ x, y int }
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if grid[i][j] == '#' {
				positions = append(positions, struct{ x, y int }{i, j})
			}
		}
	}

	// A valid tetromino must have exactly 4 '#' characters.
	if len(positions) != 4 {
		return false
	}

	// Check if all '#' characters are connected.
	visited := make(map[int]bool)
	queue := []int{0}
	visited[0] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for i := 0; i < len(positions); i++ {
			if !visited[i] && isAdjacent(positions[current], positions[i]) {
				visited[i] = true
				queue = append(queue, i)
			}
		}
	}

	// If all 4 positions are visited, they are connected.
	return len(visited) == 4
}

// isAdjacent checks if two positions are adjacent.
func isAdjacent(a, b struct{ x, y int }) bool {
	dx := a.x - b.x
	dy := a.y - b.y
	return (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
}

func main() {
	// Read the file.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [4][4]rune
	row := 0

	// Read the file line by line.
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 4 {
			fmt.Println("Invalid line length")
			return
		}
		for col, char := range line {
			if char != '#' && char != '.' {
				fmt.Println("Invalid character in line")
				return
			}
			grid[row][col] = char
		}
		row++
		if row == 4 {
			break
		}
	}

	if row != 4 {
		fmt.Println("Not enough lines in file")
		return
	}

	// Check if the grid contains a valid tetromino.
	if isValidTetromino(grid) {
		fmt.Println("Valid tetromino")
	} else {
		fmt.Println("Invalid tetromino")
	}
}
