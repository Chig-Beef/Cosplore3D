package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	// Enemies
	load_image(g, "blob1", "blobFront")
	load_image(g, "blob2", "blobRight")
	load_image(g, "blob3", "blobBack")
	load_image(g, "blob4", "blobLeft")
	load_image(g, "crewmate1", "crewmateFront")
	load_image(g, "crewmate2", "crewmateRight")
	load_image(g, "crewmate3", "crewmateBack")
	load_image(g, "crewmate4", "crewmateLeft")
	load_image(g, "crawler1", "crawlerFront")
	load_image(g, "crawler2", "crawlerRight")
	load_image(g, "crawler3", "crawlerBack")
	load_image(g, "crawler4", "crawlerLeft")

	// HUD
	load_image(g, "heart", "heart")
	load_image(g, "gun", "gun")

	// Walls
	load_image(g, "cosplorerWall", "cosplorerWall")
	load_image(g, "ankaranWall", "ankaranWall")
	load_image(g, "cosplorerComputer", "cosplorerComputer")
	load_image(g, "cosplorerReactor", "cosplorerReactor")
	load_image(g, "enikokoWall", "enikokoWall")

	// Items
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

func create_image_columns(g *Game, keys []string) {
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

func apply_image_colums_to_tiles(g *Game, l *Level) {
	for row := 0; row < len(l.data); row++ {
		for col := 0; col < len(l.data[row]); col++ {
			l.data[row][col].imageCols = g.imageColumns[get_tile_image(l.data[row][col].code)]
		}
	}
	l.fullyLoaded = true
}
