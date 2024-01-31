package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1280
	screenHeight = 640
	tileSize     = 200
)

type Game struct {
	levels  map[string]Level
	player  *Player
	enemies []Enemy
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.player.update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	g.player.draw(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}

	g.levels = load_levels(tileSize)

	camera := Camera{}
	camera.fov = 90
	camera.dov = 2_000

	g.player = &Player{
		tileSize * 3.5,
		tileSize * 3.5,
		0,
		&camera,
		"test",
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func to_radians(angle float64) float64 {
	return angle * math.Pi / 180.0
}

func bound_angle(angle *float64) {
	for *angle < 0 {
		*angle += 360
	}
	for *angle > 360 {
		*angle -= 360
	}
}
