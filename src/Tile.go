package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Tile struct {
	init      bool
	x         float64
	y         float64
	w         float64
	h         float64
	code      uint8
	imageCols *[]*ebiten.Image
}

func (t *Tile) check_hit(x, y int) bool {

	if x == int(t.x) || x == int(t.x+t.w) {
		if y >= int(t.y) && y <= int(t.y+t.h) {
			return true
		}
	}
	if y == int(t.y) || y == int(t.y+t.h) {
		if x >= int(t.x) && x <= int(t.x+t.w) {
			return true
		}
	}

	return false
}

func get_color_with_distance(c color.Color, d float64) color.Color {
	clr, ok := c.(color.RGBA)
	if !ok {
		return c
	}

	modifier := math.Sqrt(d / float64(tileSize))
	if modifier < 1 {
		modifier = 1
	}

	return color.RGBA{
		uint8(float64(clr.R) / modifier),
		uint8(float64(clr.G) / modifier),
		uint8(float64(clr.B) / modifier),
		uint8(clr.A),
	}
}

func create_image_columns(g *Game, keys []string) {
	g.imageColumns[""] = &[]*ebiten.Image{}

	for _, key := range keys {
		img := g.images[key]

		imgW, imgH := img.Size()
		images := make([]*ebiten.Image, imgW) // Since we know the size we can initialise with size
		var newImage *ebiten.Image
		for i := 0; i < imgW; i++ {
			newImage, _ = ebiten.NewImage(1, imgH, ebiten.FilterDefault)
			for j := 0; j < imgH; j++ {
				clr := img.At(i, j)
				newImage.Set(0, j, clr)
			}
			images[i] = newImage
		}

		g.imageColumns[key] = &images
	}
}
