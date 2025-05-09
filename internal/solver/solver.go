package solver

import (
	"fmt"
	"sort"
)

const (
	maxBoardSize = 10
)

func SolveTetrominos(tetrominos []*Tetromino) (string, error) {
	if len(tetrominos) == 0 {
		return "", fmt.Errorf("ERROR")
	}

	minArea := len(tetrominos) * 4
	minSize := 2
	for minSize*minSize < minArea {
		minSize++
	}

	originalOrder := make([]*Tetromino, len(tetrominos))
	copy(originalOrder, tetrominos)

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
			if board == nil {
				continue
			}
			if solution, solved := solveWithoutRotation(tetrominos, 0, board); solved {
				return boardToString(solution), nil
			}
			copy(tetrominos, originalOrder)
		}
	}

	var dimensions []struct{ width, height int }
	for height := 2; height <= maxBoardSize; height++ {
		for width := height + 1; width <= maxBoardSize; width++ {
			if width*height >= minArea {
				dimensions = append(dimensions, struct{ width, height int }{width, height})
			}
		}
	}

	sort.Slice(dimensions, func(i, j int) bool {
		ratioI := float64(dimensions[i].width) / float64(dimensions[i].height)
		ratioJ := float64(dimensions[j].width) / float64(dimensions[j].height)
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
			if board == nil {
				continue
			}
			if solution, solved := solveWithoutRotation(tetrominos, 0, board); solved {
				return boardToString(solution), nil
			}
			copy(tetrominos, originalOrder)
		}
	}

	return "", fmt.Errorf("ERROR")
}

func solveWithoutRotation(tetrominos []*Tetromino, index int, board *Board) (*Board, bool) {
	if index >= len(tetrominos) {
		return board, true
	}

	current := tetrominos[index]
	if current == nil || board == nil {
		return nil, false
	}

	maxY := board.Height - current.Height
	maxX := board.Width - current.Width

	if maxY < 0 || maxX < 0 {
		return nil, false
	}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
