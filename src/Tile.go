package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

func (t *Tile) check_point_in_tile(x, y float64) bool {
	return x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h
}

func (t *Tile) check_line_in_tile(x1, y1, x2, y2 float64) bool {
	x3, y3 := t.x, t.y
	x4, y4 := t.x, t.y+t.h

	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			if line_interset(x1, y1, x2, y2, x3, y3, x4, y4) {
				return true
			}
			x4 += t.w
			y4 -= t.h
		}
		x3 += t.w
		y3 += t.h
		x4, y4 = t.x, t.y+t.h
	}

	return false
}

/*func line_interset(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	// Calculate the slopes
	m1 := (y2 - y1) / (x2 - x1)
	m2 := (y4 - y3) / (x4 - x3)

	// Parallel (or close to apparently?)
	if math.Abs(m1-m2) < 1e-6 {
		return false
	}

	// Int as in intersection, not integer idiot
	var xInt, yInt float64
	if math.IsInf(m1, 0) {
		xInt = x1
		yInt = m2*(x1-x3) + y3
	} else if math.IsInf(m2, 0) {
		xInt = x3
		yInt = m1*(x3-x1) + y1
	} else {
		xInt = (m1*x1 - m2*x3 + y3 - y1) / (m1 - m2)
		yInt = m1*(xInt-x1) + y1
	}

	// Massive check
	return math.Min(x1, x2) <= xInt && xInt <= math.Max(x1, x2) &&
		math.Min(y1, y2) <= yInt && yInt <= math.Max(y1, y2) &&
		math.Min(x3, x4) <= xInt && xInt <= math.Max(x3, x4) &&
		math.Min(y3, y4) <= yInt && yInt <= math.Max(y3, y4)
}*/

func line_interset(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))

	if 0 <= t && t <= 1 && 0 <= u && u <= 1 {
		return true
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
			newImage = ebiten.NewImage(1, imgH)
			for j := 0; j < imgH; j++ {
				clr := img.At(i, j)
				newImage.Set(0, j, clr)
			}
			images[i] = newImage
		}

		g.imageColumns[key] = &images
	}
}
