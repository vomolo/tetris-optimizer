package functions

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	tetrisDir = "tetris_files"
	minLines  = 4
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

	if len(content) < 20 {
		return "", newValidationError("file too small to be valid")
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
		points     []Point
		minX, maxX = 3, 0
		minY, maxY = 3, 0
	)

	for y, line := range block {
		for x, char := range line {
			if char == '#' {
				hashCount++
				points = append(points, Point{X: x, Y: y})
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

	for i := range points {
		points[i].X -= minX
		points[i].Y -= minY
	}

	return &Tetromino{
		Points: points,
		Letter: 'A' + rune(blockNumber),
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
	}, nil
}

func SolveTetrominos(tetrominos []*Tetromino) (string, error) {
	minSize := 2
	for minSize*minSize < len(tetrominos)*4 {
		minSize++
	}

	sort.Slice(tetrominos, func(i, j int) bool {
		return tetrominos[i].Width*tetrominos[i].Height >
			tetrominos[j].Width*tetrominos[j].Height
	})

	for size := minSize; size <= minSize+5; size++ {
		board := NewBoard(size)
		if solution, solved := solve(tetrominos, 0, board); solved {
			return boardToString(solution), nil
		}
	}

	return "", fmt.Errorf("no solution found")
}

func solve(tetrominos []*Tetromino, index int, board *Board) (*Board, bool) {
    if index == len(tetrominos) {
        return board, true
    }

    current := tetrominos[index]
    rotations := generateRotations(current)

    for _, rot := range rotations {
        // Ensure we don't go out of bounds when placing
        maxY := board.Size - rot.Height
        maxX := board.Size - rot.Width
        
        // Skip if piece is larger than board
        if maxY < 0 || maxX < 0 {
            continue
        }

        for y := 0; y <= maxY; y++ {
            for x := 0; x <= maxX; x++ {
                if board.canPlace(rot, x, y) {
                    board.place(rot, x, y)
                    if solution, solved := solve(tetrominos, index+1, board); solved {
                        return solution, true
                    }
                    board.remove(rot, x, y)
                }
            }
        }
    }

    return nil, false
}

func generateRotations(t *Tetromino) []*Tetromino {
	rotations := []*Tetromino{t}

	for i := 0; i < 3; i++ {
		rotated := &Tetromino{
			Letter: t.Letter,
			Width:  t.Height,
			Height: t.Width,
		}

		for _, p := range rotations[len(rotations)-1].Points {
			rotated.Points = append(rotated.Points, Point{
				X: t.Height - 1 - p.Y,
				Y: p.X,
			})
		}
		rotations = append(rotations, rotated)
	}

	return rotations
}

func NewBoard(size int) *Board {
	grid := make([][]rune, size)
	for i := range grid {
		grid[i] = make([]rune, size)
	}
	return &Board{Grid: grid, Size: size}
}

func (b *Board) canPlace(t *Tetromino, x, y int) bool {
    for _, p := range t.Points {
        nx, ny := x+p.X, y+p.Y
        // Check both positive bounds and upper bounds
        if nx < 0 || ny < 0 || nx >= b.Size || ny >= b.Size || b.Grid[ny][nx] != 0 {
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
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Grid[y][x] == 0 {
				buf.WriteByte('.')
			} else {
				buf.WriteRune(b.Grid[y][x])
			}
		}
		if y < b.Size-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}