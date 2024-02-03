package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Player struct {
	x        float64
	y        float64
	angle    float64
	camera   *Camera
	curLevel string
	speed    float64
	haste    float64
}

func (p *Player) update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.angle -= p.haste
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.angle += p.haste
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.x += math.Cos(to_radians(p.angle)) * p.speed
		p.y += math.Sin(to_radians(p.angle)) * p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.x -= math.Cos(to_radians(p.angle)) * p.speed
		p.y -= math.Sin(to_radians(p.angle)) * p.speed
	}

	bound_angle(&p.angle)

	p.camera.update_props(p)
}

func (p *Player) draw(g *Game, screen *ebiten.Image) {
	p.camera.draw_world(g.levels[p.curLevel], screen)

	p.draw_hud()

	p.camera.draw_2D(g.levels[p.curLevel], screen)
}

func (p *Player) draw_hud() {

}
