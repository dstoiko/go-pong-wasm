package main

import (
	"fmt"
	"github.com/dstoiko/go-pong-wasm/pong"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"runtime"
)

type gameState byte

const (
	startState gameState = iota
	controlsState
	playState
	interState
	pauseState
	gameOverState
)

var state = startState

// Game is the structure of the game state
type Game struct {
	state    gameState
	aiMode   bool
	ball     *pong.Ball
	player1  *pong.Paddle
	player2  *pong.Paddle
	rally    int
	level    int
	maxScore int
}

const (
	initBallVelocity = 5.0
	initPaddleSpeed  = 10.0
	speedUpdateCount = 6
	speedIncrement   = 0.5
)

const (
	windowWidth   = 800
	windowHeight  = 600
	fontSize      = 30
	smallFontSize = int(fontSize / 2)
)

var (
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

var bgColor = color.Black
var objColor = color.RGBA{120, 226, 160, 255}

// NewGame creates an initializes a new game
func NewGame(aiMode bool) *Game {
	g := &Game{}
	g.init(aiMode)
	return g
}

func (g *Game) init(aiMode bool) {
	g.state = startState
	g.aiMode = aiMode
	if aiMode {
		g.maxScore = 100
	} else {
		g.maxScore = 11
	}

	g.player1 = &pong.Paddle{
		Position: pong.Position{
			X: pong.InitPaddleShift,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  pong.InitPaddleWidth,
		Height: pong.InitPaddleHeight,
		Color:  objColor,
		Up:     ebiten.KeyW,
		Down:   ebiten.KeyS,
	}
	g.player2 = &pong.Paddle{
		Position: pong.Position{
			X: windowWidth - pong.InitPaddleShift - pong.InitPaddleWidth,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  pong.InitPaddleWidth,
		Height: pong.InitPaddleHeight,
		Color:  objColor,
		Up:     ebiten.KeyO,
		Down:   ebiten.KeyK,
	}
	g.ball = &pong.Ball{
		Position: pong.Position{
			X: float32(windowWidth / 2),
			Y: float32(windowHeight / 2)},
		Radius:    pong.InitBallRadius,
		Color:     objColor,
		XVelocity: initBallVelocity,
		YVelocity: initBallVelocity,
	}
	g.level = 0
	g.ball.Img, _ = ebiten.NewImage(int(g.ball.Radius*2), int(g.ball.Radius*2), ebiten.FilterDefault)
	g.player1.Img, _ = ebiten.NewImage(g.player1.Width, g.player1.Height, ebiten.FilterDefault)
	g.player2.Img, _ = ebiten.NewImage(g.player2.Width, g.player2.Height, ebiten.FilterDefault)

	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	var dpi float64 = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(smallFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func (g *Game) reset(screen *ebiten.Image, state gameState) {
	w, _ := screen.Size()
	g.state = state
	g.rally = 0
	g.level = 0
	if state == startState {
		g.player1.Score = 0
		g.player2.Score = 0
	}
	g.player1.Position = pong.Position{
		X: pong.InitPaddleShift, Y: pong.GetCenter(screen).Y}
	g.player2.Position = pong.Position{
		X: float32(w - pong.InitPaddleShift - pong.InitPaddleWidth), Y: pong.GetCenter(screen).Y}
	g.ball.Position = pong.GetCenter(screen)
	g.ball.XVelocity = initBallVelocity
	g.ball.YVelocity = initBallVelocity
}

// Update updates the game state
func (g *Game) Update(screen *ebiten.Image) error {
	switch g.state {
	case startState:
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.state = controlsState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.aiMode = true
			g.state = playState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			g.aiMode = false
			g.state = playState
		}

	case controlsState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = startState
		}
	case playState:
		w, _ := screen.Size()

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pauseState
			break
		}

		g.player1.Update(screen)
		if g.aiMode {
			g.player2.AiUpdate(g.ball)
		} else {
			g.player2.Update(screen)
		}

		xV := g.ball.XVelocity
		g.ball.Update(g.player1, g.player2, screen)
		// rally count
		if xV*g.ball.XVelocity < 0 {
			// score up when ball touches human player's paddle
			if g.aiMode && g.ball.X < float32(w/2) {
				g.player1.Score++
			}

			g.rally++

			// spice things up
			if (g.rally)%speedUpdateCount == 0 {
				g.level++
				g.ball.XVelocity += speedIncrement
				g.ball.YVelocity += speedIncrement
				g.player1.Speed += speedIncrement
				g.player2.Speed += speedIncrement
			}
		}

		if g.ball.X < 0 {
			g.player2.Score++
			if g.aiMode {
				g.state = gameOverState
				break
			}
			g.reset(screen, interState)
		} else if g.ball.X > float32(w) {
			g.player1.Score++
			if g.aiMode {
				g.state = gameOverState
				break
			}
			g.reset(screen, interState)
		}

		if g.player1.Score == g.maxScore || g.player2.Score == g.maxScore {
			g.state = gameOverState
		}

	case interState, pauseState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = playState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset(screen, startState)
		}

	case gameOverState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset(screen, startState)
		}
	}

	g.Draw(screen)

	return nil
}

func (g *Game) drawCaption(screen *ebiten.Image) {
	w, h := screen.Size()
	msg := []string{}
	switch g.state {
	case playState, interState, pauseState:
		msg = append(msg, "Press SPACE key to take a break (not too long though)")
	case controlsState:
		msg = append(msg, "Press SPACE to go back to main menu")
	}
	for i, l := range msg {
		x := (w - len(l)*smallFontSize) / 2
		text.Draw(screen, l, smallArcadeFont, x, h-4+(i-2)*smallFontSize, objColor)
	}
}

func (g *Game) drawBigText(screen *ebiten.Image) {
	w, _ := screen.Size()
	var texts []string
	switch g.state {
	case startState:
		texts = []string{
			"",
			"PONG",
			"",
			"C -> CONTROLS",
			"V -> VS GAME",
			"A -> AI GAME",
		}
	case controlsState:
		texts = []string{
			"",
			"PLAYER 1:",
			"W -> UP",
			"S -> DOWN",
			"",
			"PLAYER 2:",
			"O -> UP",
			"K -> DOWN",
		}
	case interState:
		texts = []string{
			"",
			"",
			"SPACE -> RESUME",
			"R     -> RESET",
		}
	case pauseState:
		texts = []string{
			"",
			"PAUSED",
			"",
			"SPACE -> RESUME",
			"R     -> RESET",
		}
	case gameOverState:
		texts = []string{
			"",
			"GAME OVER!",
		}
		if g.player1.Score == g.maxScore {
			texts = append(texts, "PLAYER 1 WINS")
		} else if g.player2.Score == g.maxScore && !g.aiMode {
			texts = append(texts, "PLAYER 2 WINS")
		} else {
			texts = append(texts, "AI WINS")
		}
		texts = append(texts, "SPACE -> RESET")
	}
	for i, l := range texts {
		x := (w - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, objColor)
	}
}

// Draw updates the game screen elements drawn
func (g *Game) Draw(screen *ebiten.Image) error {
	screen.Fill(bgColor)

	g.drawCaption(screen)
	g.drawBigText(screen)

	if g.state != controlsState {
		g.player1.Draw(screen, arcadeFont, false)
		g.player2.Draw(screen, arcadeFont, g.aiMode)
		g.ball.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

// Layout sets the screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func main() {
	// On browsers, let's use fullscreen so that this is playable on any browsers.
	// It is planned to ignore the given 'scale' apply fullscreen automatically on browsers (#571).
	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	}
	ai := true
	g := NewGame(ai)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
