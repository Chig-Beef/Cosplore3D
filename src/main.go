package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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
	ctx             *audio.Context
	musicPlayer     map[string]*AudioPlayer
	musicPlayerCh   map[string]chan *AudioPlayer
	errCh           chan error
	audio           map[string]*AudioPlayer
	curAudio        string
	soundEffects    map[string]*AudioPlayer
	curSoundEffects []string
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
		for i := 0; i < len(g.levels[g.player.curLevel].bosses); i++ {
			g.levels[g.player.curLevel].bosses[i].update(g, g.levels[g.player.curLevel].get_solid_tiles())
		}

		g.player.update(g)
		g.kull_enemies()
		g.kull_bosses()
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
	g.levels[g.player.curLevel].draw_items_enemies_and_bosses(screen, g.player.camera)
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

	if err := g.init_audio(); err != nil {
		log.Fatal(err)
	}

	g.menu = Menu{
		[]*Button{
			{0, screenWidth / 15 * 2, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, "Play", start_game},
			{0, screenWidth / 15 * 3, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, "Settings", start_game},
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
