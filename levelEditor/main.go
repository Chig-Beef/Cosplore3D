package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1280
	screenHeight = 640
	tileSize     = 50.0
)

type Game struct {
	images       map[string]*ebiten.Image
	fonts        map[string]font.Face
	curMousePos  [2]int
	prevMousePos [2]int

	curCodeSelection uint8

	level Level
}

type Level struct {
	name       string
	data       [][]uint8
	floorColor color.RGBA
	skyColor   color.RGBA
}

func (g *Game) Update() error {

	x, y := ebiten.CursorPosition()

	if g.curMousePos == [2]int{} {
		g.curMousePos = [2]int{1, 1}
	} else {
		g.curMousePos = [2]int{x, y}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		col := int(math.Floor(float64(g.curMousePos[0]) / tileSize))
		row := int(math.Floor(float64(g.curMousePos[1]-160) / tileSize))

		if row < len(g.level.data) && row >= 0 {
			if col < len(g.level.data[row]) && col >= 0 {
				g.level.data[row][col] = g.curCodeSelection
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		g.curCodeSelection++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		col := int(math.Floor(float64(g.curMousePos[0]) / tileSize))
		row := int(math.Floor(float64(g.curMousePos[1]-160) / tileSize))

		if row < len(g.level.data) && row >= 0 {
			if col < len(g.level.data[row]) && col >= 0 {
				g.level.data[row][col] = 0
			}
		}
	}

	g.prevMousePos = g.curMousePos
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for row := 0; row < len(g.level.data); row++ {
		for col := 0; col < len(g.level.data[row]); col++ {
			draw_relevant_image(screen, g, g.level.data[row][col], row, col)
		}
	}
}

func draw_relevant_image(screen *ebiten.Image, g *Game, code uint8, row, col int) {
	switch code {
	case 0:
		ebitenutil.DrawRect(screen, float64(col*tileSize), 160+float64(row)*tileSize, tileSize, tileSize, color.Black)
	case 1:
		img := g.images["cosplorerWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col)*tileSize, float64(row)*tileSize+160)

		screen.DrawImage(img, &op)
	case 2:
		img := g.images["ankaranWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col)*tileSize, float64(row)*tileSize+160)

		screen.DrawImage(img, &op)
	case 7:
		img := g.images["cosmium"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col)*tileSize, float64(row)*tileSize+160)

		screen.DrawImage(img, &op)
	case 8:
		img := g.images["blobFront"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col)*tileSize, float64(row)*tileSize+160)

		screen.DrawImage(img, &op)
	case 9:
		img := g.images["blobFront"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col)*tileSize, float64(row)*tileSize+160)

		screen.DrawImage(img, &op)
	default:
		ebitenutil.DrawRect(screen, float64(col*tileSize), 160+float64(row)*tileSize, tileSize, tileSize, color.RGBA{255, 0, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	load_image(g, "blob1", "blobFront")
	load_image(g, "cosplorerWall", "cosplorerWall")
	load_image(g, "ankaranWall", "ankaranWall")
	load_image(g, "cosmium", "cosmium")
	load_image(g, "ammo", "ammo")
}

func load_image(g *Game, fName string, mName string) {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/" + fName + ".png")
	if err != nil {
		log.Fatal(err)
	}
	g.images[mName] = img
}

func blank_level(width, height int) [][]uint8 {
	output := make([][]uint8, height)
	for row := 0; row < height; row++ {
		output[row] = make([]uint8, width)
	}
	output[5][5] = 1
	return output
}

func main() {

	g := &Game{}

	g.curCodeSelection = 1

	g.load_images()

	g.level = Level{
		"unknown",
		blank_level(32, 32),
		color.RGBA{},
		color.RGBA{},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CosEditor")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
