package main

import (
	"image/color"
	"math"
)

type Tile struct {
	x     float64
	y     float64
	w     float64
	h     float64
	code  uint8
	color color.Color
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
