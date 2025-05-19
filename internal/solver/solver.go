package solver

import (
	"fmt"
	"math"
)

const (
	maxBoardSize = 20
)

func SolveTetrominos(tetrominos []*Tetromino) (string, error) {
	if len(tetrominos) == 0 {
		return "", NewValidationError("no tetrominos provided")
	}

	minArea := len(tetrominos) * 4

	originalOrder := make([]*Tetromino, len(tetrominos))
	copy(originalOrder, tetrominos)

	sortStrategies := []func(i, j int) bool{
		func(i, j int) bool {
			return tetrominos[i].Width*tetrominos[i].Height > tetrominos[j].Width*tetrominos[j].Height
		},
		func(i, j int) bool { return tetrominos[i].Height > tetrominos[j].Height },
		func(i, j int) bool { return tetrominos[i].Width > tetrominos[j].Width },
	}

	type Dim struct{ W, H int }

	var dimensions []Dim
	for area := minArea; area <= maxBoardSize*maxBoardSize; area++ {
		for w := 1; w <= maxBoardSize; w++ {
			if area%w == 0 {
				h := area / w
				if h <= maxBoardSize {
					dimensions = append(dimensions, Dim{w, h})
				}
			}
		}
	}

	sort.Slice(dimensions, func(i, j int) bool {
		ai := dimensions[i].W * dimensions[i].H
		aj := dimensions[j].W * dimensions[j].H
		if ai != aj {
			return ai < aj
		}
		if dimensions[i].W != dimensions[j].W {
			return dimensions[i].W < dimensions[j].W
		}
		return dimensions[i].H < dimensions[j].H
	})

	for _, dim := range dimensions {
		for _, sortFn := range sortStrategies {
			sort.Slice(tetrominos, sortFn)
			board := NewBoard(dim.W, dim.H)
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

	return "", fmt.Errorf("no optimized square solution found")
}
