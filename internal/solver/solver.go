package solver

import (
	"fmt"
	"math"
	"sort"
)

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidationError(message string) error {
	return &ValidationError{message: message}
}

func SolveTetrominos(tetrominos []*Tetromino) (string, error) {
	if len(tetrominos) == 0 {
		return "", NewValidationError("no tetrominos provided")
	}

	// Assign unique letters to each tetromino
	for i, t := range tetrominos {
		t.Letter = rune('A' + i)
	}

	// First try optimized solution for repetitive tetrominos
	if solution, err := tryOptimizedSquareRepetitiveSolution(tetrominos); err == nil {
		return solution, nil
	}

	// Fall back to general solver
	return generalSquareSolver(tetrominos)
}

func tryOptimizedSquareRepetitiveSolution(tetrominos []*Tetromino) (string, error) {
	groups := groupRepetitiveTetrominos(tetrominos)
	if len(groups) != 1 || len(groups[0].tetrominos) < 5 {
		return "", fmt.Errorf("not a repetitive case")
	}

	t := groups[0].tetrominos[0]
	n := len(tetrominos)
	totalBlocks := n * 4
	minSize := int(math.Ceil(math.Sqrt(float64(totalBlocks))))

	// Try to find the smallest square that can fit all pieces
	for size := minSize; size <= minSize+5; size++ {
		// Calculate how many pieces fit in rows and columns
		piecesPerRow := size / t.Width
		piecesPerCol := size / t.Height
		totalPieces := piecesPerRow * piecesPerCol

		if totalPieces < n {
			continue // Not enough space
		}

		board := NewBoard(size)
		if board == nil {
			continue
		}

		success := true
		placed := 0
		for row := 0; row < piecesPerCol && placed < n; row++ {
			for col := 0; col < piecesPerRow && placed < n; col++ {
				x := col * t.Width
				y := row * t.Height
				if !board.CanPlace(tetrominos[placed], x, y) {
					success = false
					break
				}
				board.Place(tetrominos[placed], x, y)
				placed++
			}
			if !success {
				break
			}
		}

		if success && placed == n {
			return board.String(), nil
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

func normalizeTetromino(t *Tetromino) []Point {
	minX, minY := t.Points[0].X, t.Points[0].Y
	for _, p := range t.Points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	normalized := make([]Point, len(t.Points))
	for i, p := range t.Points {
		normalized[i] = Point{p.X - minX, p.Y - minY}
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Y == normalized[j].Y {
			return normalized[i].X < normalized[j].X
		}
		return normalized[i].Y < normalized[j].Y
	})

	return normalized
}

func calculateComplexity(t *Tetromino) int {
	complexity := 0
	for _, p := range t.Points {
		neighbors := 0
		for _, d := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			np := Point{p.X + d.X, p.Y + d.Y}
			if containsPoint(t.Points, np) {
				neighbors++
			}
		}
		complexity += (4 - neighbors)
	}
	return complexity
}

func containsPoint(points []Point, p Point) bool {
	for _, point := range points {
		if point == p {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}