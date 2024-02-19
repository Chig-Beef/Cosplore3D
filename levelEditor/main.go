package main

import (
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

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
	tileSize     = 50.0
)

type Game struct {
	images       map[string]*ebiten.Image
	fonts        map[string]font.Face
	curMousePos  [2]int
	prevMousePos [2]int

	offset      [2]int
	cameraSpeed int

	textBoxes [6]TextBox

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
		if 600 <= g.curMousePos[0] && g.curMousePos[0] <= 650 {
			if 10 <= g.curMousePos[1] && g.curMousePos[1] <= 60 {
				save(g, g.level.data)
				return nil
			}
		}

		for i := 0; i < len(g.textBoxes); i++ {
			if g.textBoxes[i].check_click(g) {
				return nil
			}
		}

		col := int(math.Floor(float64(g.curMousePos[0]+g.offset[0]) / tileSize))
		row := int(math.Floor(float64(g.curMousePos[1]-160+g.offset[1]) / tileSize))

		if row < len(g.level.data) && row >= 0 {
			if col < len(g.level.data[row]) && col >= 0 {
				g.level.data[row][col] = g.curCodeSelection
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.curCodeSelection++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.curCodeSelection--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		g.curCodeSelection -= 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.curCodeSelection += 10
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		col := int(math.Floor(float64(g.curMousePos[0]+g.offset[0]) / tileSize))
		row := int(math.Floor(float64(g.curMousePos[1]-160+g.offset[1]) / tileSize))

		if row < len(g.level.data) && row >= 0 {
			if col < len(g.level.data[row]) && col >= 0 {
				g.level.data[row][col] = 0
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.offset[0] -= g.cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.offset[0] += g.cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.offset[1] -= g.cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.offset[1] += g.cameraSpeed
	}

	for i := 0; i < len(g.textBoxes); i++ {
		g.textBoxes[i].update()
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

	ebitenutil.DrawRect(screen, 600, 10, 50, 50, color.White)

	for i := 0; i < 6; i++ {
		g.textBoxes[i].draw(screen, g)
	}
}

func draw_relevant_image(screen *ebiten.Image, g *Game, code uint8, row, col int) {
	switch code {
	case 0:
		ebitenutil.DrawRect(screen, float64(col*tileSize-g.offset[0]), 160+float64(row*tileSize-g.offset[1]), tileSize, tileSize, color.Black)

	// Walls
	case 1:
		img := g.images["ankaranWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 2:
		img := g.images["cosplorerWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 3:
		img := g.images["computer"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 4:
		img := g.images["reactor"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 5:
		img := g.images["enikokoWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 6:
		img := g.images["schmeltoolWall"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)

	// Markers
	case 10:
		img := g.images["heart"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)

	// Items
	case 30:
		img := g.images["ammo"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 31:
		img := g.images["cosmium"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)

	// Enemies
	case 40:
		img := g.images["blobFront"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 41:
		img := g.images["crewmateLeft"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 42:
		img := g.images["crawlerFront"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	case 43:
		img := g.images["beastFront"]

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(tileSize/float64(img.Bounds().Dx()), tileSize/float64(img.Bounds().Dy()))
		op.GeoM.Translate(float64(col*tileSize-g.offset[0]), float64(row*tileSize-g.offset[1])+160)

		screen.DrawImage(img, &op)
	default:
		ebitenutil.DrawRect(screen, float64(col*tileSize-g.offset[0]), 160+float64(row*tileSize-g.offset[1]), tileSize, tileSize, color.RGBA{255, 0, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	load_image(g, "blob1", "blobFront")
	load_image(g, "crewmate2", "crewmateLeft")
	load_image(g, "crawler1", "crawlerFront")
	load_image(g, "beast1", "beastFront")
	load_image(g, "cosplorerWall", "cosplorerWall")
	load_image(g, "ankaranWall", "ankaranWall")
	load_image(g, "enikokoWall", "enikokoWall")
	load_image(g, "schmeltoolWall", "schmeltoolWall")
	load_image(g, "cosmium", "cosmium")
	load_image(g, "ammo", "ammo")
	load_image(g, "heart", "heart")
	load_image(g, "cosplorerComputer", "computer")
	load_image(g, "cosplorerReactor", "reactor")
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
	return output
}

func save(g *Game, d [][]uint8) {
	output := ""
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			output += strconv.Itoa(int(g.textBoxes[i*3+j].value)) + ","
		}
		output = output[:len(output)-1] + "|"
	}
	output = output[:len(output)-1] + "\r\n"

	for row := 0; row < len(d); row++ {
		for col := 0; col < len(d[row]); col++ {
			output += strconv.Itoa(int(d[row][col]))
			output += ","
		}
		// Get rid of extra comma
		output = output[:len(output)-1] + "\r\n"
	}
	// Get rid of extra newline
	output = output[:len(output)-2]

	file, err := os.Create("temp.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write([]byte(output))
	if err != nil {
		log.Fatal(err)
	}
}

func load(game *Game) Level {
	name := "unknown"
	floorColor := color.RGBA{128, 128, 128, 255}
	skyColor := color.RGBA{0, 0, 0, 255}

	rawData, err := os.ReadFile("temp.txt")
	if err != nil {
		log.Println("Could not find or read temp.txt")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	splitData := strings.Split(string(rawData), "\r\n")

	colors := strings.Split(splitData[0], "|")

	rawFloorColor := strings.Split(colors[0], ",")
	//fmt.Println(rawData)

	if len(rawFloorColor) != 3 {
		log.Println("Need 3 color arguments in floor color")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	r, err := strconv.Atoi(rawFloorColor[0])
	if err != nil {
		log.Println("Need integers for colors floor red")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	g, err := strconv.Atoi(rawFloorColor[1])
	if err != nil {
		log.Println("Need integers for colors floor green")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	b, err := strconv.Atoi(rawFloorColor[2])
	if err != nil {
		log.Println("Need integers for colors floor blue")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	floorColor.R = uint8(r)
	floorColor.G = uint8(g)
	floorColor.B = uint8(b)

	game.textBoxes[0].set_value(r)
	game.textBoxes[1].set_value(g)
	game.textBoxes[2].set_value(b)

	rawSkyColor := strings.Split(colors[1], ",")

	if len(rawSkyColor) != 3 {
		log.Println("Need 3 color arguments in sky color")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	r, err = strconv.Atoi(rawSkyColor[0])
	if err != nil {
		log.Println("Need integers for colors sky red")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	g, err = strconv.Atoi(rawSkyColor[1])
	if err != nil {
		log.Println("Need integers for colors sky green")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	b, err = strconv.Atoi(rawSkyColor[2])
	if err != nil {
		log.Println("Need integers for colors sky blue")
		return Level{name, blank_level(32, 32), floorColor, skyColor}
	}

	skyColor.R = uint8(r)
	skyColor.G = uint8(g)
	skyColor.B = uint8(b)

	game.textBoxes[3].set_value(r)
	game.textBoxes[4].set_value(g)
	game.textBoxes[5].set_value(b)

	data := [][]uint8{}

	rows := splitData[1:]
	for y := 0; y < len(rows); y++ {
		data = append(data, []uint8{})
		row := strings.Split(rows[y], ",")
		for x := 0; x < len(row); x++ {
			n, err := strconv.Atoi(row[x])
			if err != nil {
				log.Println("Need integers for data", x, y)
				return Level{name, blank_level(32, 32), floorColor, skyColor}
			}
			data[y] = append(data[y], uint8(n))
		}
	}

	return Level{name, data, floorColor, skyColor}
}

func (g *Game) load_fonts() {
	g.fonts = make(map[string]font.Face)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 20
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight/8 - 20,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["colors"] = mplusNormalFont
}

func (g *Game) create_textboxes() {
	g.textBoxes = [6]TextBox{}
	for i := 0; i < 6; i++ {
		g.textBoxes[i] = TextBox{
			"0",
			0,
			float64(i * 62),
			5,
			60,
			20,
			color.RGBA{64, 64, 64, 255},
			color.White,
			false,
		}
	}
}

func main() {

	g := &Game{}

	g.curCodeSelection = 1

	g.load_images()
	g.load_fonts()
	g.create_textboxes()

	g.level = load(g)

	g.cameraSpeed = 5

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CosEditor")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
