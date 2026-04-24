package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"
	"math"

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

var (
	switchWidth    = 160.0
	switchHeight   = 36.0
	switchX        = float64(screenWidth)/2 - switchWidth/2
	switchY        = 370.0
	switchTrackOff = color.RGBA{180, 160, 100, 255}
	switchTrackOn  = color.RGBA{100, 160, 100, 255}
	switchThumbCol = color.RGBA{255, 255, 255, 230}
)

func isMouseOverButton(x, y, width, height float64) bool {
	mx, my := ebiten.CursorPosition()
	return float64(mx) >= x && float64(mx) <= x+width &&
		float64(my) >= y && float64(my) <= y+height
}

func drawSwitch(screen *ebiten.Image, aiMode bool) {
	trackColor := switchTrackOff
	if aiMode {
		trackColor = switchTrackOn
	}
	vector.FillRect(screen,
		float32(switchX+switchHeight/2), float32(switchY),
		float32(switchWidth-switchHeight), float32(switchHeight),
		trackColor, true)
	vector.FillCircle(screen,
		float32(switchX+switchHeight/2), float32(switchY+switchHeight/2),
		float32(switchHeight/2), trackColor, true)
	vector.FillCircle(screen,
		float32(switchX+switchWidth-switchHeight/2), float32(switchY+switchHeight/2),
		float32(switchHeight/2), trackColor, true)
	thumbCX := switchX + switchHeight/2
	if aiMode {
		thumbCX = switchX + switchWidth - switchHeight/2
	}
	vector.FillCircle(screen,
		float32(thumbCX), float32(switchY+switchHeight/2),
		float32(switchHeight/2-3), switchThumbCol, true)
}

func titleScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})
	if g.winner != 0 {
		if g.winner == 1 {
			g.text.Draw(screen, "X wins!", 380, 280)
		} else if g.winner == 2 {
			g.text.Draw(screen, "O wins!", 380, 280)
		} else if g.winner == 3 {
			g.text.Draw(screen, "It's a draw!", 380, 280)
		}
	}
	if titleImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(titleX, titleY)
		screen.DrawImage(titleImage, op)
	}
	buttonColor := grey
	if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) {
		buttonColor = darkGrey
	}
	if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) ||
		isMouseOverButton(switchX, switchY, switchWidth, switchHeight) {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	vector.FillRect(screen,
		float32(buttonX), float32(buttonY),
		float32(buttonWidth), float32(buttonHeight),
		buttonColor, true)
	if playTextImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.6, 0.6)
		op.GeoM.Translate(playTextX+30, playTextY+15)
		screen.DrawImage(playTextImage, op)
	}
	drawSwitch(screen, g.aiMode)
	g.text.Draw(screen, "2P", int(switchX-10)-30, int(switchY+switchHeight/2)+8)
	g.text.Draw(screen, "vs AI", int(switchX+switchWidth+10), int(switchY+switchHeight/2)+8)
}

func drawGrid(screen *ebiten.Image, g *Game) {
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

func gameScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})
	drawGrid(screen, g)
}

func gameAIScene(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{242, 240, 239, 255})
	drawGrid(screen, g)
}

type Game struct {
	state        uint8
	turn         uint8
	player       uint8
	board        [3][3]uint8
	winner       uint8
	text         *etxt.Renderer
	aiMode       bool
	endGameTimer uint8
}

func (g *Game) init() error {
	g.turn = 0
	g.player = 0
	g.board = [3][3]uint8{}
	g.winner = 0
	g.aiMode = false
	g.state = 0
	var err error
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
		if isMouseOverButton(switchX, switchY, switchWidth, switchHeight) {
			g.aiMode = !g.aiMode
			return
		}
		if isMouseOverButton(buttonX, buttonY, float64(buttonWidth), float64(buttonHeight)) {
			g.turn = 0
			g.player = 0
			g.winner = 0
			g.board = [3][3]uint8{}
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			if g.aiMode {
				g.state = 2
			} else {
				g.state = 1
			}
			return
		}
	}
}

func checkWin(board [3][3]uint8) int {
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
					switch checkWin(g.board) {
					case 1:
						log.Println("X wins!")
						g.state = 3
						g.winner = 1
						g.endGameTimer = 0
						return
					case 2:
						log.Println("O wins!")
						g.state = 3
						g.winner = 2
						g.endGameTimer = 0
						return
					case 0:
						log.Println("It's a draw!")
						g.state = 3
						g.winner = 3
						g.endGameTimer = 0
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

func gameAIUpdate(g *Game) {
	if g.turn != g.player {
		aiMove := getBestMove(g.board)
		if aiMove != -1 {
			g.board[aiMove/3][aiMove%3] = 2
		}
		g.turn = g.player
		switch checkWin(g.board) {
		case 1:
			log.Println("X wins!")
			g.state = 3
			g.winner = 1
			g.endGameTimer = 0
			return
		case 2:
			log.Println("O wins!")
			g.state = 3
			g.winner = 2
			g.endGameTimer = 0
			return
		case 0:
			log.Println("It's a draw!")
			g.state = 3
			g.winner = 3
			g.endGameTimer = 0
			return
		}
		return
	}
	overCell := false
	for i := 0; i < 9; i++ {
		x, y, w, h := gridBoxes[i][0], gridBoxes[i][1], gridBoxes[i][2], gridBoxes[i][3]
		if isMouseOverButton(x, y, w, h) {
			overCell = true
			break
		}
	}
	if overCell {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for i := 0; i < 9; i++ {
			x, y, w, h := gridBoxes[i][0], gridBoxes[i][1], gridBoxes[i][2], gridBoxes[i][3]
			if isMouseOverButton(x, y, w, h) && g.board[i%3][i/3] == 0 {
				g.board[i%3][i/3] = 1
				g.turn = 1
				switch checkWin(g.board) {
				case 1:
					log.Println("X wins!")
					g.state = 3
					g.winner = 1
					g.endGameTimer = 0
					return
				case 2:
					log.Println("O wins!")
					g.state = 3
					g.winner = 2
					g.endGameTimer = 0
					return
				case 0:
					log.Println("It's a draw!")
					g.state = 3
					g.winner = 3
					g.endGameTimer = 0
					return
				}
				break
			}
		}
	}
}

func minimax(board [3][3]uint8, isMaximizing bool) int {
	result := checkWin(board)
	if result == 2 {
		return 10
	} else if result == 1 {
		return -10
	} else if result == 0 {
		return 0
	}
	if isMaximizing {
		bestScore := -1000
		for i := 0; i < 9; i++ {
			row := i / 3
			col := i % 3
			if board[row][col] == 0 {
				board[row][col] = 2
				score := minimax(board, false)
				board[row][col] = 0
				if score > bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := 0; i < 9; i++ {
			row := i / 3
			col := i % 3
			if board[row][col] == 0 {
				board[row][col] = 1
				score := minimax(board, true)
				board[row][col] = 0
				if score < bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	}
}

func getBestMove(board [3][3]uint8) int {
	bestScore := -1000
	bestMove := -1
	for i := 0; i < 9; i++ {
		row := i / 3
		col := i % 3
		if board[row][col] == 0 {
			board[row][col] = 2
			score := minimax(board, false)
			board[row][col] = 0
			if score > bestScore {
				bestScore = score
				bestMove = i
			}
		}
	}
	return bestMove
}

func endGameUpdate(g *Game) {
	g.endGameTimer++
	if g.endGameTimer >= 120 {
		g.state = 0
	}
}

func (g *Game) Update() error {
	switch g.state {
	case 0:
		titleUpdate(g)
	case 1:
		gameUpdate(g)
	case 2:
		gameAIUpdate(g)
	case 3:
		endGameUpdate(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case 0:
		titleScene(screen, g)
	case 1:
		gameScene(screen, g)
	case 2:
		gameAIScene(screen, g)
	case 3:
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
