package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

	// Audio
	musicPlayer   *AudioPlayer
	musicPlayerCh chan *AudioPlayer
	errCh         chan error
}

func (g *Game) Update() error {
	//fmt.Println(ebiten.ActualTPS())
	if !g.hasLoadedImageColumns {
		create_image_columns(g, []string{
			"cosplorerWall",
			"ankaranWall",
			"cosplorerComputer",
		})
		apply_image_colums_to_tiles(g, g.levels[g.player.curLevel])
		return nil
	}

	select {
	case p := <-g.musicPlayerCh:
		g.musicPlayer = p
	case err := <-g.errCh:
		return err
	default:
	}

	if g.musicPlayer != nil {
		if err := g.musicPlayer.update(g); err != nil {
			return err
		}
	}

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

	for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
		g.levels[g.player.curLevel].enemies[i].update(g, g.levels[g.player.curLevel].get_solid_tiles())
	}

	g.player.update(g)

	g.kull_enemies()

	for i := 0; i < len(g.levels[g.player.curLevel].progressors); i++ {
		g.levels[g.player.curLevel].progressors[i].check_collide(g)
	}

	g.prevMousePos = g.curMousePos
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.hasLoadedImageColumns {
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

func (g *Game) load_fonts() {
	g.fonts = make(map[string]font.Face)

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

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
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
	}

	g.player = &Player{
		tileSize * 3.5,
		tileSize * 3.5,
		0,
		&camera,
		"cosplorer",
		10,
		3,
		100,
		weapon,
		[]string{},
	}

	g.levels = load_levels(g, tileSize)

	g.player.x = g.levels[g.player.curLevel].playerStartPos[0]
	g.player.y = g.levels[g.player.curLevel].playerStartPos[1]

	pre_init_image_columns(g, []string{
		"cosplorerWall",
		"ankaranWall",
		"cosplorerComputer",
	})

	g.init_audio()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
