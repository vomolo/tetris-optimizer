package solver

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
