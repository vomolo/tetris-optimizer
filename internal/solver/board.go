package solver

import "bytes"

// Board represents the Tetris game board.
type Board struct {
	Grid   [][]rune
	Size   int // Square board, so width = height
	Placed int
}

// NewBoard creates a new square board of given size.
func NewBoard(size int) *Board {
	if size <= 0 {
		return nil
	}
	grid := make([][]rune, size)
	for i := range grid {
		grid[i] = make([]rune, size)
	}
	return &Board{Grid: grid, Size: size}
}

// CanPlace checks if a tetromino can be placed at position (x, y).
func (b *Board) CanPlace(t *Tetromino, x, y int) bool {
	for _, p := range t.Points {
		nx, ny := x+p.X, y+p.Y
		if nx < 0 || ny < 0 || nx >= b.Size || ny >= b.Size || b.Grid[ny][nx] != 0 {
			return false
		}
	}
	return true
}

// Place places a tetromino at position (x, y).
func (b *Board) Place(t *Tetromino, x, y int) {
	for _, p := range t.Points {
		b.Grid[y+p.Y][x+p.X] = t.Letter
	}
	b.Placed++
}

// Remove removes a tetromino from position (x, y).
func (b *Board) Remove(t *Tetromino, x, y int) {
	for _, p := range t.Points {
		b.Grid[y+p.Y][x+p.X] = 0
	}
	b.Placed--
}

// String converts the board to a string representation.
func (b *Board) String() string {
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
