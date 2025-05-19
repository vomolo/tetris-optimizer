package solver

import (
	"fmt"
	"math"
	"sort"
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

func generalSquareSolver(tetrominos []*Tetromino) (string, error) {
	// Sort tetrominos by size and complexity
	sortedTetrominos := make([]*Tetromino, len(tetrominos))
	copy(sortedTetrominos, tetrominos)
	sort.Slice(sortedTetrominos, func(i, j int) bool {
		areaI := sortedTetrominos[i].Width * sortedTetrominos[i].Height
		areaJ := sortedTetrominos[j].Width * sortedTetrominos[j].Height
		if areaI != areaJ {
			return areaI > areaJ
		}
		return calculateComplexity(sortedTetrominos[i]) > calculateComplexity(sortedTetrominos[j])
	})

	totalBlocks := len(tetrominos) * 4
	minSize := int(math.Ceil(math.Sqrt(float64(totalBlocks))))

	// Try solving with increasing square board sizes
	for size := minSize; size <= minSize+5; size++ {
		board := NewBoard(size)
		if board == nil {
			continue
		}
		if solve(board, sortedTetrominos, 0) {
			return board.String(), nil
		}
	}
	return "", fmt.Errorf("no square solution found")
}

func solve(board *Board, tetrominos []*Tetromino, index int) bool {
	if index == len(tetrominos) {
		return true
	}

	t := tetrominos[index]
	for y := 0; y <= board.Size-t.Height; y++ {
		for x := 0; x <= board.Size-t.Width; x++ {
			if !board.CanPlace(t, x, y) {
				continue
			}

			board.Place(t, x, y)
			if solve(board, tetrominos, index+1) {
				return true
			}
			board.Remove(t, x, y)
		}
	}
	return false
}

type tetrominoGroup struct {
	tetrominos []*Tetromino
	points     []Point
}

func groupRepetitiveTetrominos(tetrominos []*Tetromino) []tetrominoGroup {
	groups := make([]tetrominoGroup, 0)
	used := make(map[int]bool)

	for i, t1 := range tetrominos {
		if used[i] {
			continue
		}
		group := tetrominoGroup{tetrominos: []*Tetromino{t1}, points: t1.Points}
		for j := i + 1; j < len(tetrominos); j++ {
			if !used[j] && areTetrominosEqual(t1, tetrominos[j]) {
				group.tetrominos = append(group.tetrominos, tetrominos[j])
				used[j] = true
			}
		}
		groups = append(groups, group)
	}
	return groups
}

func areTetrominosEqual(t1, t2 *Tetromino) bool {
	if len(t1.Points) != len(t2.Points) {
		return false
	}
	norm1 := normalizeTetromino(t1)
	norm2 := normalizeTetromino(t2)
	for i := range norm1 {
		if norm1[i] != norm2[i] {
			return false
		}
	}
	return true
}
