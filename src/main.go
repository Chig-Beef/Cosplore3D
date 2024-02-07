package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1280
	screenHeight = 640
	tileSize     = 200
)

type Game struct {
	hasLoadedLevels bool
	levels          map[string]Level
	player          *Player
	enemies         []Enemy
	images          map[string]*ebiten.Image
	imageColumns    map[string]*[]*ebiten.Image
	fonts           map[string]font.Face
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !g.hasLoadedLevels {
		g.levels = load_levels(g, tileSize)
	}
	g.player.update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.levels[g.player.curLevel].draw_floor_sky(screen)
	g.player.draw(g, screen)
	for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
		g.levels[g.player.curLevel].enemies[i].draw(screen, *g.player.camera)
	}
	g.player.draw_hud(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) load_fonts() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight/8 - 20,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["ammo"] = mplusNormalFont
}

func (g *Game) load_images() {
	load_image(g, "blob1", "blob1")
	load_image(g, "heart", "heart")
	load_image(g, "gun", "gun")
	load_image(g, "cosplorerWall", "cosplorerWall")
}

func load_image(g *Game, fName string, mName string) {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/"+fName+".png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	g.images[mName] = img
}

func main() {
	g := &Game{}

	g.images = make(map[string]*ebiten.Image)
	g.fonts = make(map[string]font.Face)
	g.imageColumns = make(map[string]*[]*ebiten.Image)

	g.load_images()
	g.load_fonts()

	camera := Camera{
		0,
		0,
		0,
		90,
		20 * tileSize,
	}

	g.player = &Player{
		tileSize * 3.5,
		tileSize * 3.5,
		0,
		&camera,
		"test",
		7,
		3,
		100,
		Weapon{},
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
