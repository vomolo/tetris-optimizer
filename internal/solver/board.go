package solver

import "bytes"

type Board struct {
	Grid   [][]rune
	Width  int
	Height int
	Placed int
}

func NewBoard(width, height int) *Board {
	if width <= 0 || height <= 0 {
		return &Board{
			Grid:   [][]rune{},
			Width:  0,
			Height: 0,
		}
	}

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