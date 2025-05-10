package solver

import "fmt"

type Point struct {
	X, Y int
}

type Tetromino struct {
	Points []Point
	Letter rune
	Width  int
	Height int
}

func validateAndCreateTetromino(block [][]byte, blockNumber int) (*Tetromino, error) {
	if len(block) != 4 {
		return nil, fmt.Errorf("ERROR")
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
			return nil, fmt.Errorf("ERROR")
		}

		for x, char := range line {
			if char == '#' {
				if hashCount >= 4 {
					return nil, fmt.Errorf("ERROR")
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
				return nil, fmt.Errorf("ERROR")
			}
		}
	}

	if hashCount != 4 {
		return nil, fmt.Errorf("ERROR")
	}

	for i := range points {
		points[i].X -= minX
		points[i].Y -= minY
	}

	if !isValidTetromino(points) {
		return nil, fmt.Errorf("ERROR")
	}

	return &Tetromino{
		Points: points[:],
		Letter: 'A' + rune(blockNumber),
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
	}, nil
}

func isValidTetromino(points [4]Point) bool {
	// Check for duplicates (optional, depending on requirements)
	for i := range 4 {
		for j := i + 1; j < 4; j++ {
			if points[i].X == points[j].X && points[i].Y == points[j].Y {
				return false // duplicate points
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

	// Check if all points are reachable from the first point (single connected component)
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
