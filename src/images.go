package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	// Enemies
	g.load_image("blob1", "blobFront")
	g.load_image("blob2", "blobRight")
	g.load_image("blob3", "blobBack")
	g.load_image("blob4", "blobLeft")
	g.load_image("crewmate1", "crewmateFront")
	g.load_image("crewmate2", "crewmateRight")
	g.load_image("crewmate3", "crewmateBack")
	g.load_image("crewmate4", "crewmateLeft")
	g.load_image("crawler1", "crawlerFront")
	g.load_image("crawler2", "crawlerRight")
	g.load_image("crawler3", "crawlerBack")
	g.load_image("crawler4", "crawlerLeft")

	// HUD
	g.load_image("heart", "heart")
	g.load_image("gun", "gun")
	g.load_image("gunFire", "gunFire")

	// Walls
	g.load_image("cosplorerWall", "cosplorerWall")
	g.load_image("ankaranWall", "ankaranWall")
	g.load_image("cosplorerComputer", "cosplorerComputer")
	g.load_image("cosplorerReactor", "cosplorerReactor")
	g.load_image("cosplorerReactorEmpty", "cosplorerReactorEmpty")
	g.load_image("enikokoWall", "enikokoWall")
	g.load_image("schmeltoolWall", "schmeltoolWall")

	// Items
	g.load_image("cosmium", "cosmium")
	g.load_image("ammo", "ammo")

	// Other
	g.load_image("planet", "planet")
}

func (g *Game) load_image(fName string, mName string) {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/" + fName + ".png")
	if err != nil {
		log.Fatal(err)
	}
	g.images[mName] = img
}

func (g *Game) create_image_columns(keys []string) {
	g.imageColumns[""] = &[]*ebiten.Image{}

	for _, key := range keys {
		img := g.images[key]

		imgW, imgH := img.Size()
		images := make([]*ebiten.Image, imgW) // Since we know the size we can initialise with size
		var newImage *ebiten.Image
		for i := 0; i < imgW; i++ {
			newImage = ebiten.NewImage(1, imgH)
			for j := 0; j < imgH; j++ {
				clr := img.At(i, j)
				newImage.Set(0, j, clr)
			}
			images[i] = newImage
		}

		g.imageColumns[key] = &images
	}

	g.hasLoadedImageColumns = true
}

func (g *Game) apply_image_columns_to_tiles(l *Level) {
	for row := 0; row < len(l.data); row++ {
		for col := 0; col < len(l.data[row]); col++ {
			l.data[row][col].imageCols = g.imageColumns[get_tile_image(l.data[row][col].code)]
		}
	}
	l.fullyLoaded = true
}
