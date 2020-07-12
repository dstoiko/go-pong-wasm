package pong

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"strconv"
)

// Paddle is a pong paddle
type Paddle struct {
	Position
	Score        int
	Speed        float32
	Width        int
	Height       int
	Color        color.Color
	Up           ebiten.Key
	Down         ebiten.Key
	Img          *ebiten.Image
	pressed      keysPressed
	scorePrinted scorePrinted
}

const (
	InitPaddleWidth  = 20
	InitPaddleHeight = 100
	InitPaddleShift  = 50
)

type keysPressed struct {
	up   bool
	down bool
}

type scorePrinted struct {
	score   int
	printed bool
	x       int
	y       int
}

func (p *Paddle) Update(screen *ebiten.Image) {
	_, h := screen.Size()

	if inpututil.IsKeyJustPressed(p.Up) {
		p.pressed.down = false
		p.pressed.up = true
	} else if inpututil.IsKeyJustReleased(p.Up) || !ebiten.IsKeyPressed(p.Up) {
		p.pressed.up = false
	}
	if inpututil.IsKeyJustPressed(p.Down) {
		p.pressed.up = false
		p.pressed.down = true
	} else if inpututil.IsKeyJustReleased(p.Down) || !ebiten.IsKeyPressed(p.Down) {
		p.pressed.down = false
	}

	if p.pressed.up {
		p.Y -= p.Speed
	} else if p.pressed.down {
		p.Y += p.Speed
	}

	if p.Y-float32(p.Height/2) < 0 {
		p.Y = float32(1 + p.Height/2)
	} else if p.Y+float32(p.Height/2) > float32(h) {
		p.Y = float32(h - p.Height/2 - 1)
	}
}

func (p *Paddle) AiUpdate(b *Ball) {
	// unbeatable haha
	p.Y = b.Y
}

func (p *Paddle) Draw(screen *ebiten.Image, scoreFont font.Face, ai bool) {
	// draw player's paddle
	pOpts := &ebiten.DrawImageOptions{}
	pOpts.GeoM.Translate(float64(p.X), float64(p.Y-float32(p.Height/2)))
	p.Img.Fill(p.Color)
	screen.DrawImage(p.Img, pOpts)

	// draw player's score if needed
	if !ai {
		if p.scorePrinted.score != p.Score && p.scorePrinted.printed {
			p.scorePrinted.printed = false
		}
		if p.scorePrinted.score == 0 && !p.scorePrinted.printed {
			p.scorePrinted.x = int(p.X + (GetCenter(screen).X-p.X)/2)
			p.scorePrinted.y = int(2 * 30)
		}
		if (p.scorePrinted.score == 0 || p.scorePrinted.score != p.Score) && !p.scorePrinted.printed {
			p.scorePrinted.score = p.Score
			p.scorePrinted.printed = true
		}
		s := strconv.Itoa(p.scorePrinted.score)
		text.Draw(screen, s, scoreFont, p.scorePrinted.x, p.scorePrinted.y, p.Color)
	}
}
