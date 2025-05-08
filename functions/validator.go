package functions

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	tetrisDir    = "tetris_files"
	minLines     = 4
	maxBoardSize = 10
)

func Validate(filename string) (string, error) {
	fullPath := filename
	if filepath.Dir(filename) != tetrisDir {
		fullPath = filepath.Join(tetrisDir, filename)
	}

	if err := validateStructure(fullPath); err != nil {
		return "", err
	}
	return validateAndSolveContent(fullPath)
}

func validateStructure(fullPath string) error {
	if filepath.Ext(fullPath) != ".txt" {
		return newValidationError("file must have .txt extension")
	}

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return newValidationError("file '%s' does not exist in %s directory",
				filepath.Base(fullPath), tetrisDir)
		}
		return newValidationError("file access error: %v", err)
	}
	return nil
}

func validateAndSolveContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", newValidationError("failed to read file: %v", err)
	}

	if len(content) < 16 {
		return "", newValidationError("ERROR")
	}

	lines := bytes.Split(content, []byte{'\n'})
	var (
		lineCount    int
		blockLines   [4][]byte
		blockIndex   int
		hasContent   bool
		blockCounter int
		tetrominos   []*Tetromino
	)

	for _, line := range lines {
		lineCount++
		trimmed := bytes.TrimSpace(line)

		for _, char := range line {
			if char != '#' && char != '.' && char != '\n' && char != '\r' && char != ' ' && char != '\t' {
				return "", newValidationError("invalid character '%c' in line %d", char, lineCount)
			}
		}

		if lineCount%5 == 0 {
			if len(trimmed) > 0 {
				return "", newValidationError("line %d must be empty", lineCount)
			}

			tetromino, err := validateAndCreateTetromino(blockLines[:], blockCounter)
			if err != nil {
				return "", err
			}
			tetrominos = append(tetrominos, tetromino)
			blockIndex = 0
			blockCounter++
			continue
		}

		if len(trimmed) == 0 {
			return "", newValidationError("line %d cannot be empty", lineCount)
		}
		if len(trimmed) != 4 {
			return "", newValidationError("line %d must have exactly 4 characters", lineCount)
		}

		if blockIndex >= 4 {
			return "", newValidationError("block %d has too many lines", blockCounter+1)
		}
		blockLines[blockIndex] = trimmed
		blockIndex++
		hasContent = true
	}

	if blockIndex > 0 {
		tetromino, err := validateAndCreateTetromino(blockLines[:blockIndex], blockCounter)
		if err != nil {
			return "", err
		}
		tetrominos = append(tetrominos, tetromino)
	}

	if !hasContent {
		return "", newValidationError("file is empty")
	}
	if lineCount < minLines {
		return "", newValidationError("file must have at least %d lines", minLines)
	}

	return SolveTetrominos(tetrominos)
}

func validateAndCreateTetromino(block [][]byte, blockNumber int) (*Tetromino, error) {
	if len(block) != 4 {
		return nil, newValidationError("block must have 4 lines (found %d)", len(block))
	}

	var (
		hashCount  int
		points     [4]Point // Fixed size array since we know it's exactly 4 points
		minX, maxX = 3, 0
		minY, maxY = 3, 0
		pointIdx   int
	)

	// First pass: count hashes and record positions
	for y, line := range block {
		for x, char := range line {
			if char == '#' {
				if hashCount >= 4 {
					return nil, newValidationError("block has more than 4 '#' characters")
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
		return nil, newValidationError("block must have exactly 4 '#' (found %d)", hashCount)
	}

	// Create normalized points (relative to minX, minY)
	for i := range points {
		points[i].X -= minX
		points[i].Y -= minY
	}

	// Validate the tetromino shape
	if !isValidTetromino(points) {
		return nil, newValidationError("invalid tetromino shape at block %d", blockNumber)
	}

	return &Tetromino{
		Points: points[:],
		Letter: 'A' + rune(blockNumber),
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
	}, nil
}

func isValidTetromino(points [4]Point) bool {
	// Precompute all possible neighbor offsets
	var connected [4]bool

	// Check each point's neighbors
	for i := 0; i < 4; i++ {
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

	// Count connected points
	connectionCount := 0
	for _, c := range connected {
		if c {
			connectionCount++
		}
	}

	// For a valid tetromino, we need at least 3 connections
	// (linear shape has 3, others have more)
	return connectionCount >= 3
}

func SolveTetrominos(tetrominos []*Tetromino) (string, error) {
	minArea := len(tetrominos) * 4
	minSize := 2
	for minSize*minSize < minArea {
		minSize++
	}

	// Cache the original order
	originalOrder := make([]*Tetromino, len(tetrominos))
	copy(originalOrder, tetrominos)

	// Try square boards first with different sorting strategies
	sortStrategies := []func(i, j int) bool{
		func(i, j int) bool {
			return tetrominos[i].Width*tetrominos[i].Height > tetrominos[j].Width*tetrominos[j].Height
		},
		func(i, j int) bool { return tetrominos[i].Height > tetrominos[j].Height },
		func(i, j int) bool { return tetrominos[i].Width > tetrominos[j].Width },
	}

	for size := minSize; size <= min(maxBoardSize, minSize+5); size++ {
		for _, sortFn := range sortStrategies {
			sort.Slice(tetrominos, sortFn)
			board := NewBoard(size, size)
			if solution, solved := solveWithoutRotation(tetrominos, 0, board); solved {
				return boardToString(solution), nil
			}
			// Restore original order for next attempt
			copy(tetrominos, originalOrder)
		}
	}

	// If no square solution found, try rectangular boards
	var dimensions []struct{ width, height int }
	for height := 2; height <= maxBoardSize; height++ {
		for width := height + 1; width <= maxBoardSize; width++ {
			if width*height >= minArea {
				dimensions = append(dimensions, struct{ width, height int }{width, height})
			}
		}
	}

	// Sort rectangular boards by how close they are to square
	sort.Slice(dimensions, func(i, j int) bool {
		ratioI := float64(dimensions[i].width) / float64(dimensions[i].height)
		ratioJ := float64(dimensions[j].width) / float64(dimensions[j].height)
		// Prefer boards closer to square
		diffI := ratioI - 1.0
		if diffI < 0 {
			diffI = -diffI
		}
		diffJ := ratioJ - 1.0
		if diffJ < 0 {
			diffJ = -diffJ
		}
		if diffI != diffJ {
			return diffI < diffJ
		}
		return dimensions[i].width*dimensions[i].height < dimensions[j].width*dimensions[j].height
	})

	for _, dim := range dimensions {
		for _, sortFn := range sortStrategies {
			sort.Slice(tetrominos, sortFn)
			board := NewBoard(dim.width, dim.height)
			if solution, solved := solveWithoutRotation(tetrominos, 0, board); solved {
				return boardToString(solution), nil
			}
			copy(tetrominos, originalOrder)
		}
	}

	return "", fmt.Errorf("no solution found")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveWithoutRotation(tetrominos []*Tetromino, index int, board *Board) (*Board, bool) {
	if index == len(tetrominos) {
		return board, true
	}

	current := tetrominos[index]
	maxY := board.Height - current.Height
	maxX := board.Width - current.Width

	if maxY < 0 || maxX < 0 {
		return nil, false
	}

	// Try placing the piece in all possible positions
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if board.canPlace(current, x, y) {
				board.place(current, x, y)
				if solution, solved := solveWithoutRotation(tetrominos, index+1, board); solved {
					return solution, true
				}
				board.remove(current, x, y)
			}
		}
	}

	return nil, false
}

func NewBoard(width, height int) *Board {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
	}
	return &Board{Grid: grid, Width: width, Height: height}
}

func (b *Board) canPlace(t *Tetromino, x, y int) bool {
	for _, p := range t.Points {
		nx, ny := x+p.X, y+p.Y
		if nx < 0 || ny < 0 || nx >= b.Width || ny >= b.Height || b.Grid[ny][nx] != 0 {
			return false
		}
	}
	return true
}

func (b *Board) place(t *Tetromino, x, y int) {
	for _, p := range t.Points {
		b.Grid[y+p.Y][x+p.X] = t.Letter
	}
	b.Placed++
}

func (b *Board) remove(t *Tetromino, x, y int) {
	for _, p := range t.Points {
		b.Grid[y+p.Y][x+p.X] = 0
	}
	b.Placed--
}

func boardToString(b *Board) string {
	var buf bytes.Buffer
	for y := range b.Height {
		for x := range b.Width {
			if b.Grid[y][x] == 0 {
				buf.WriteByte('.')
			} else {
				buf.WriteRune(b.Grid[y][x])
			}
		}
		if y < b.Height-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
