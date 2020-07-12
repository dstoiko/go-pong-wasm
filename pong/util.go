package pong

import (
	"github.com/hajimehoshi/ebiten"
)

// Position is a set of coordinates in 2-D plan
type Position struct {
	X, Y float32
}

// GetCenter returns the center position on screen
func GetCenter(screen *ebiten.Image) Position {
	w, h := screen.Size()
	return Position{
		X: float32(w / 2),
		Y: float32(h / 2),
	}
}
