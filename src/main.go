package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	images  map[string]ebiten.Image
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.player.update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DrawRect(screen, 0, screenHeight/2, screenWidth, screenHeight/2, color.Gray{32})
	g.player.draw(g, screen)
	for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
		g.levels[g.player.curLevel].enemies[i].draw(screen, *g.player.camera)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}

	g.images = make(map[string]ebiten.Image)

	img, _, err := ebitenutil.NewImageFromFile("assets/images/blob1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	g.images["blob1"] = *img

	g.levels = load_levels(g, tileSize)

	camera := Camera{}
	camera.fov = 90
	camera.dov = 2_000

	g.player = &Player{
		tileSize * 3.5,
		tileSize * 3.5,
		0,
		&camera,
		"test",
		5,
		2,
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

func to_degrees(angle float64) float64 {
	return angle * 180.0 / math.Pi
}

func bound_angle(angle *float64) float64 {
	for *angle < 0 {
		*angle += 360
	}
	for *angle > 360 {
		*angle -= 360
	}
	return *angle
}
