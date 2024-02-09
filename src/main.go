package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	hasLoadedLevels bool
	levels          map[string]*Level
	player          *Player
	images          map[string]*ebiten.Image
	imageColumns    map[string]*[]*ebiten.Image
	fonts           map[string]font.Face
	curMousePos     [2]int
	prevMousePos    [2]int

	// Audio
	musicPlayer   *AudioPlayer
	musicPlayerCh chan *AudioPlayer
	errCh         chan error
}

func (g *Game) Update() error {
	//fmt.Println(ebiten.ActualTPS())

	if !g.hasLoadedLevels {
		g.levels = load_levels(g, tileSize)
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
		g.levels[g.player.curLevel].enemies[i].update(g)
	}

	g.player.update(g)

	g.kull_enemies()

	g.prevMousePos = g.curMousePos
	return nil
}

func (g *Game) kull_enemies() {
	var e *Enemy
	aliveEnemies := []int{}
	for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
		e = g.levels[g.player.curLevel].enemies[i]
		if e.health > 0 {
			aliveEnemies = append(aliveEnemies, i)
		}
	}
	newEnemies := make([]*Enemy, len(aliveEnemies))
	for i := 0; i < len(aliveEnemies); i++ {
		newEnemies[i] = g.levels[g.player.curLevel].enemies[aliveEnemies[i]]
	}
	g.levels[g.player.curLevel].enemies = newEnemies
}

func (g *Game) Draw(screen *ebiten.Image) {
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

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	load_image(g, "blob1", "blob1")
	load_image(g, "heart", "heart")
	load_image(g, "gun", "gun")
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

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	g := &Game{}

	g.imageColumns = make(map[string]*[]*ebiten.Image)

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
		"ankaran",
		10,
		3,
		100,
		weapon,
	}

	g.init_audio()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cosplore3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
