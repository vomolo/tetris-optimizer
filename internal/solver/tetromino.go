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

func validateAndCreateTetromino(block [][]byte, blockNumber int) (*Tetromino) {
	if len(block) != 4 {
		fmt.Println("ERROR")
		return nil
	}

	var (
		hashCount  int
		points     [4]Point
		minX, maxX = 3, 0
		minY, maxY = 3, 0
		pointIdx   int
	)

	for y, line := range block {
		for x, char := range line {
			if char == '#' {
				if hashCount >= 4 {
					fmt.Println("ERROR")
					return nil
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
			}
		}
	}

	if hashCount != 4 {
		fmt.Println("ERROR")
		return nil
	}

	for i := range points {
		points[i].X -= minX
		points[i].Y -= minY
	}

	if !isValidTetromino(points) {
		fmt.Println("ERROR")
		return nil
	}

	return &Tetromino{
		Points: points[:],
		Letter: 'A' + rune(blockNumber),
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
	}, nil
}

func isValidTetromino(points [4]Point) bool {
	var connected [4]bool

	for i := range 4 {
		for j := i + 1; j < 4; j++ {
			dx := points[i].X - points[j].X
			dy := points[i].Y - points[j].Y
			if (dx == 1 && dy == 0) || (dx == -1 && dy == 0) ||
				(dx == 0 && dy == 1) || (dx == 0 && dy == -1) {
				connected[i] = true
				connected[j] = true
			}
		}
	}

	connectionCount := 0
	for _, c := range connected {
		if c {
			connectionCount++
		}
	}

	return connectionCount >= 3
}