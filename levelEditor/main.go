package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1280
	screenHeight = 640
)

type Game struct {
	images       map[string]*ebiten.Image
	fonts        map[string]font.Face
	curMousePos  [2]int
	prevMousePos [2]int
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

	g.prevMousePos = g.curMousePos
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
