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
	settings              Setting
	sensitivity           float64
	difficulty            float64
	ending                Ending

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
	} else if g.player.curLevel == "settings" {
		g.settings.update(g)
	} else if g.player.curLevel == "ending" {
		g.ending.update()
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
	switch g.player.curLevel {
	case "menu":
		g.menu.draw(screen, g)
	case "settings":
		g.settings.draw(screen, g)
	case "ending":
		g.ending.draw(screen)
	default:
		g.levels[g.player.curLevel].draw_floor_sky(screen)
		g.player.draw(g, screen)
		g.levels[g.player.curLevel].draw_items_enemies_and_bosses(screen, g.player.camera)
		g.player.draw_hud(g, screen)
	}
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
		0,
		&camera,
		"menu",
		10,
		3,
		100,
		weapon,
		[]InvItem{},
		0,
	}

	g.levels = load_levels(g, tileSize)

	if err := g.init_audio(); err != nil {
		log.Fatal(err)
	}

	g.menu = Menu{
		[]*Button{
			{0, screenWidth / 15 * 2, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, "Play", start_game},
			{0, screenWidth / 15 * 3, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, "Settings", open_settings},
		},
	}

	g.settings = Setting{
		[]*Button{
			{0, screenWidth / 15 * 2, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, "Menu", open_menu},
		},
		map[string]*TextBox{
			"fov": {"90", "FOV: ", 0, screenWidth / 15 * 3, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, false},
			"dov": {"30", "DOV: ", 0, screenWidth / 15 * 4, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, false},
			"sen": {"100", "SENSITIVITY: ", 0, screenWidth / 15 * 5, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, false},
			"dif": {"100", "DIFFICULTY: ", 0, screenWidth / 15 * 6, screenWidth / 5, screenHeight / 14, color.RGBA{32, 32, 32, 255}, color.White, false},
		},
	}

	g.ending = Ending{
		[]*ebiten.Image{
			g.images["ending0"],
			g.images["ending1"],
		},
		0,
		600,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
