package solver

import "fmt"

// Point represents a coordinate in a tetromino.
type Point struct {
	X, Y int
}

// Tetromino represents a Tetris piece with 4 blocks.
type Tetromino struct {
	Points []Point
	Letter rune
	Width  int
	Height int
}

// ValidateAndCreateTetromino creates a tetromino from a 4x4 block.
func ValidateAndCreateTetromino(block [][]byte, blockNumber int) (*Tetromino, error) {
	if len(block) != 4 {
		return nil, fmt.Errorf("tetromino must have 4 rows")
	}

	var (
		hashCount  int
		points     [4]Point
		minX, maxX = 3, 0
		minY, maxY = 3, 0
		pointIdx   int
	)

	for y, line := range block {
		if len(line) != 4 {
			return nil, fmt.Errorf("tetromino row %d must have 4 columns", y)
		}
		for x, char := range line {
			if char == '#' {
				if hashCount >= 4 {
					return nil, fmt.Errorf("tetromino has too many blocks")
				}
				points[pointIdx] = Point{X: x, Y: y}
				pointIdx++
				hashCount++
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			} else if char != '.' {
				return nil, fmt.Errorf("invalid character '%c' in tetromino", char)
			}
		}
	}

	if hashCount != 4 {
		return nil, fmt.Errorf("tetromino must have exactly 4 blocks, got %d", hashCount)
	}

	// Normalize points to top-left
	for i := range points {
		points[i].X -= minX
		points[i].Y -= minY
	}

	if !isValidTetromino(points) {
		return nil, fmt.Errorf("tetromino is not connected")
	}

	return &Tetromino{
		Points: points[:],
		Letter: 'A' + rune(blockNumber),
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
	}, nil
}

// isValidTetromino checks if the points form a valid, connected tetromino.
func isValidTetromino(points [4]Point) bool {
	// Check for duplicates
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if points[i].X == points[j].X && points[i].Y == points[j].Y {
				return false
			}
		}
	}

	// Build adjacency list
	adj := make([][]int, 4)
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			dx := points[i].X - points[j].X
			dy := points[i].Y - points[j].Y
			if (dx == 1 && dy == 0) || (dx == -1 && dy == 0) ||
				(dx == 0 && dy == 1) || (dx == 0 && dy == -1) {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	// BFS to check connectivity
	visited := make([]bool, 4)
	queue := []int{0}
	visited[0] = true
	count := 1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range adj[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				count++
			}
		}
	}

	return count == 4
}