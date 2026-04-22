package main

import (
	"bytes"
	"image/color"
	"log"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed assets/title.png
var titlePNG []byte
var titleImage *ebiten.Image
var titleX, titleY, baseTitleY float64

//go:embed assets/playText.png
var playTextPNG []byte
var playTextImage *ebiten.Image
var playTextX, playTextY float64
var (
	buttonWidth  = 120
	buttonHeight = 50
	screenWidth  = 854
	screenHeight = 480
)

var (
	buttonX  = float64((screenWidth - buttonWidth) / 2)
	buttonY  = 300.0
	grey     = color.RGBA{200, 180, 120, 255}
	darkGrey = color.RGBA{160, 140, 80, 255}
)

func isMouseOverButton() bool {
	mx, my := ebiten.CursorPosition()
	return float64(mx) >= buttonX && float64(mx) <= buttonX+float64(buttonWidth) &&
		float64(my) >= buttonY && float64(my) <= buttonY+float64(buttonHeight)
}

func titleScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})

	if titleImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			titleX, titleY)

		screen.DrawImage(titleImage, op)
	}

	buttonColor := grey
	if isMouseOverButton() {
		buttonColor = darkGrey
	}
	if isMouseOverButton() {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	vector.FillRect(screen, float32(buttonX), float32(buttonY), float32(buttonWidth), float32(buttonHeight), buttonColor, true)

	if playTextImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.6, 0.6)
		op.GeoM.Translate(playTextX+30, playTextY+15)

		screen.DrawImage(playTextImage, op)
	}
}

type Game struct {
	state uint8
}

func (g *Game) init() error {
	var err error
	g.state = 0
	titleImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(titlePNG))
	w, h := titleImage.Bounds().Dx(), titleImage.Bounds().Dy()
	titleX, baseTitleY = (854-float64(w))/2, (480-float64(h)-200)/2
	titleY = baseTitleY
	if err != nil {
		return err
	}

	playTextImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(playTextPNG))
	if err != nil {
		return err
	}
	pw, ph := playTextImage.Bounds().Dx(), playTextImage.Bounds().Dy()
	playTextX, playTextY = buttonX+float64((buttonWidth-pw)/2), buttonY+float64((buttonHeight-ph)/2)

	return nil
}

func (g *Game) Update() error {
	offset := math.Sin(float64(ebiten.Tick()) * 0.05)
	titleY = baseTitleY + offset*8

	if g.state == 0 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isMouseOverButton() {
			g.state = 1
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case 0:
		titleScene(screen, g)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 854, 480
}

func main() {
	game := &Game{}

	if err := game.init(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(854, 480)
	ebiten.SetWindowTitle("Triple T")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
