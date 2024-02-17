package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1280
	screenHeight = 640
	tileSize     = 200
)

type Game struct {
	hasLoadedImageColumns bool
	levels                map[string]*Level
	player                *Player
	images                map[string]*ebiten.Image
	imageColumns          map[string]*[]*ebiten.Image
	fonts                 map[string]font.Face
	curMousePos           [2]int
	prevMousePos          [2]int
	menu                  Menu

	// Audio
	musicPlayer   *AudioPlayer
	musicPlayerCh chan *AudioPlayer
	errCh         chan error
}

func (g *Game) Update() error {
	//fmt.Println(ebiten.ActualTPS())
	g.update_audio()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if ebiten.CursorMode() == ebiten.CursorModeCaptured {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	}

	x, y := ebiten.CursorPosition()

	if g.curMousePos == [2]int{} {
		g.curMousePos = [2]int{1, 1}
	} else {
		g.curMousePos = [2]int{x, y}
	}

	if g.player.curLevel == "menu" {
		g.menu.update(g)
	} else {
		for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
			g.levels[g.player.curLevel].enemies[i].update(g, g.levels[g.player.curLevel].get_solid_tiles())
		}

		g.player.update(g)
		g.kull_enemies()
		g.levels[g.player.curLevel].update_progressors_and_triggers(g)
	}

	g.prevMousePos = g.curMousePos
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.player.curLevel == "menu" {
		g.menu.draw(screen, g)
		return
	}
	g.levels[g.player.curLevel].draw_floor_sky(screen)
	g.player.draw(g, screen)
	g.levels[g.player.curLevel].draw_items_and_enemies(screen, g.player.camera)
	g.player.draw_hud(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	g := &Game{}

	g.imageColumns = make(map[string]*[]*ebiten.Image)
	g.hasLoadedImageColumns = false

	g.load_images()
	g.load_fonts()

	camera := Camera{
		0,
		0,
		0,
		90,
		30 * tileSize,
	}

	weapon := Weapon{
		10,
		60,
		10,
		10,
		10,
		g.images["gun"],
		g.images["gunFire"],
		0,
	}

	g.player = &Player{
		tileSize * 3.5,
		tileSize * 3.5,
		0,
		&camera,
		"menu",
		10,
		3,
		100,
		weapon,
		[]InvItem{},
	}

	g.levels = load_levels(g, tileSize)

	g.init_audio()

	g.menu = Menu{
		[]*Button{
			{float64(screenWidth/2 - 150), 200, 300, 100, color.RGBA{32, 32, 32, 255}, color.White, "Play", start_game},
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
