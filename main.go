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
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/title.png
var titlePNG []byte
var titleImage *ebiten.Image
var titleX, titleY, baseTitleY float64

//go:embed assets/x.png
var xPng []byte
var xImage *ebiten.Image

//go:embed assets/o.png
var oPng []byte
var oImage *ebiten.Image

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

var gridBoxes [9][4]float64

func isMouseOverButton(x, y, width, height float64) bool {
	mx, my := ebiten.CursorPosition()
	return float64(mx) >= x && float64(mx) <= x+width &&
		float64(my) >= y && float64(my) <= y+height
}

func titleScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})
	if g.winner != 0 {
		if g.winner == 1 {
			g.text.Draw(screen, "X wins!", 380, 380)
		} else if g.winner == 2 {
			g.text.Draw(screen, "O wins!", 380, 380)
		} else if g.winner == 3 {
			g.text.Draw(screen, "It's a draw!", 380, 380)
		}
	}
	if titleImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			titleX, titleY)

		screen.DrawImage(titleImage, op)
	}

	buttonColor := grey
	if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) {
		buttonColor = darkGrey
	}
	if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) {
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
func gameScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})
	gridSize := 300.0
	cellSize := gridSize / 3
	xStart := (float64(screenWidth) - gridSize) / 2
	yStart := (float64(screenHeight) - gridSize) / 2
	lineColor := color.RGBA{0, 0, 0, 255}
	for i := 0; i <= 3; i++ {
		x := xStart + float64(i)*cellSize
		vector.StrokeLine(screen, float32(x), float32(yStart), float32(x), float32(yStart+gridSize), 2, lineColor, false)
	}
	for j := 0; j <= 3; j++ {
		y := yStart + float64(j)*cellSize
		vector.StrokeLine(screen, float32(xStart), float32(y), float32(xStart+gridSize), float32(y), 2, lineColor, false)
	}
	k := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gridBoxes[k] = [4]float64{xStart + float64(i)*cellSize, yStart + float64(j)*cellSize, cellSize, cellSize}
			k++
		}
	}
	for i := 0; i < 9; i++ {
		x, y, w, h := gridBoxes[i][0], gridBoxes[i][1], gridBoxes[i][2], gridBoxes[i][3]
		if g.board[i%3][i/3] == 1 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.2, 0.2)
			op.GeoM.Translate(x+cellSize*0.1, y+cellSize*0.1)
			screen.DrawImage(xImage, op)
		}
		if g.board[i%3][i/3] == 2 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.2, 0.2)
			op.GeoM.Translate(x+cellSize*0.1, y+cellSize*0.1)
			screen.DrawImage(oImage, op)
		}
		if isMouseOverButton(x, y, w, h) {
			vector.FillRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{70, 70, 255, 102}, true)

		}

	}
}

type Game struct {
	state  uint8
	turn   uint8
	player uint8
	board  [3][3]uint8
	winner uint8
	text   *etxt.Renderer
}

func (g *Game) init() error {
	g.turn = 0
	g.player = 0
	g.board = [3][3]uint8{}
	g.winner = 0
	var err error
	g.state = 0
	titleImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(titlePNG))
	xImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(xPng))
	oImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(oPng))
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

	g.text, err = loadTextRenderer()
	if err != nil {
		return err
	}

	return nil
}

func loadTextRenderer() (*etxt.Renderer, error) {
	font, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	renderer := etxt.NewRenderer()
	renderer.SetFont(font)
	renderer.SetSize(24)
	renderer.SetColor(color.Black)

	return renderer, nil
}
func titleUpdate(g *Game) {

	offset := math.Sin(float64(ebiten.Tick()) * 0.05)
	titleY = baseTitleY + offset*8

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) {
			g.state = 1
			g.turn = 0
			g.player = 0
			g.winner = 0
			g.board = [3][3]uint8{}
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	}
}
func checkWin(board [3][3]uint8, g *Game) int {

	for i := 0; i < 3; i++ {
		if board[i][0] == 1 && board[i][1] == 1 && board[i][2] == 1 {
			return 1
		}
		if board[i][0] == 2 && board[i][1] == 2 && board[i][2] == 2 {
			return 2
		}
	}

	for j := 0; j < 3; j++ {
		if board[0][j] == 1 && board[1][j] == 1 && board[2][j] == 1 {
			return 1
		}
		if board[0][j] == 2 && board[1][j] == 2 && board[2][j] == 2 {
			return 2
		}
	}
	if board[0][0] == 1 && board[1][1] == 1 && board[2][2] == 1 {
		return 1
	}
	if board[0][0] == 2 && board[1][1] == 2 && board[2][2] == 2 {
		return 2
	}
	if board[0][2] == 1 && board[1][1] == 1 && board[2][0] == 1 {
		return 1
	}
	if board[0][2] == 2 && board[1][1] == 2 && board[2][0] == 2 {
		return 2
	}
	full := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				full = false
			}
		}
	}
	if full {
		log.Println("It's a draw!")
		return 0
	}
	return -1
}

func gameUpdate(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for i := 0; i < 9; i++ {
			x, y, w, h := gridBoxes[i][0], gridBoxes[i][1], gridBoxes[i][2], gridBoxes[i][3]
			if isMouseOverButton(x, y, w, h) {
				if g.turn == g.player && g.board[i%3][i/3] == 0 {
					g.board[i%3][i/3] = g.player + 1

					switch checkWin(g.board, g) {
					case 1:
						log.Println("X wins!")
						g.state = 0
						g.winner = 1
						return
					case 2:
						log.Println("O wins!")
						g.state = 0
						g.winner = 2
						return
					case 0:
						log.Println("It's a draw!")
						g.state = 0
						g.winner = 3
						return
					}

				}

				g.turn = (g.turn + 1) % 2
				g.player = (g.player + 1) % 2
				break
			}
		}
	}
}

func (g *Game) Update() error {
	switch g.state {
	case 0:
		titleUpdate(g)
	case 1:
		gameUpdate(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case 0:
		titleScene(screen, g)

	case 1:
		gameScene(screen, g)
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
