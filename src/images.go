package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	var imageData [][]string
	rawImageData, err := os.ReadFile("assets/images/images.json")
	if err != nil {
		log.Fatal("couldn't load images")
	}

	err = json.Unmarshal(rawImageData, &imageData)
	if err != nil {
		log.Fatal("failed to load `./assets/images/images.json`, file may have been tampered with, reinstall advised")
	}

	for i := 0; i < len(imageData); i++ {
		fName := imageData[i][0]
		mName := imageData[i][1]
		g.load_image(fName, mName)
	}
}

func (g *Game) load_image(fName string, mName string) {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/" + fName)
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
