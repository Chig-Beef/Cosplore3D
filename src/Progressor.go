package main

import (
	"slices"
)

type Progressor struct {
	x         float64
	y         float64
	w         float64
	h         float64
	check     Check
	levelName string
}

type Check func(*Game) bool

func (p *Progressor) check_collide(g *Game) {
	x := g.player.x
	y := g.player.y

	if p.x <= x && x <= p.x+p.w {
		if p.y <= y && y <= p.y+p.h {
			if p.check(g) {
				g.open_level(p.levelName)
			}
		}
	}
}

func has_cosmium(g *Game) bool {
	return slices.Contains(g.player.inventory, "Cosmium")
}

func has_no_cosmium(g *Game) bool {
	return !slices.Contains(g.player.inventory, "Cosmium")
}
