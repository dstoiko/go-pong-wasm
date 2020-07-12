package pong

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

// Ball is a pong ball
type Ball struct {
	Position
	Radius    float32
	XVelocity float32
	YVelocity float32
	Color     color.Color
	Img       *ebiten.Image
}

const (
	InitBallRadius = 10.0
)

func setBallPixels(c color.Color, ballImg *ebiten.Image) {
	// TODO: set pixels for round effect
	ballImg.Fill(c)
}

func (b *Ball) Update(leftPaddle *Paddle, rightPaddle *Paddle, screen *ebiten.Image) {
	_, h := screen.Size()
	b.X += b.XVelocity
	b.Y += b.YVelocity

	// bounce off edges when getting to top/bottom
	if b.Y-b.Radius > float32(h) {
		b.YVelocity = -b.YVelocity
		b.Y = float32(h) - b.Radius
	} else if b.Y+b.Radius < 0 {
		b.YVelocity = -b.YVelocity
		b.Y = b.Radius
	}

	// bounce off paddles
	if b.X-b.Radius < leftPaddle.X+float32(leftPaddle.Width/2) &&
		b.Y > leftPaddle.Y-float32(leftPaddle.Height/2) &&
		b.Y < leftPaddle.Y+float32(leftPaddle.Height/2) {
		b.XVelocity = -b.XVelocity
		b.X = leftPaddle.X + float32(leftPaddle.Width/2) + b.Radius
	} else if b.X+b.Radius > rightPaddle.X-float32(rightPaddle.Width/2) &&
		b.Y > rightPaddle.Y-float32(rightPaddle.Height/2) &&
		b.Y < rightPaddle.Y+float32(rightPaddle.Height/2) {
		b.XVelocity = -b.XVelocity
		b.X = rightPaddle.X - float32(rightPaddle.Width/2) - b.Radius
	}
}

func (b *Ball) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(b.X), float64(b.Y))
	setBallPixels(b.Color, b.Img)
	screen.DrawImage(b.Img, opts)
}
