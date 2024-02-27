package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Weapon struct {
	name             string
	damage           float64
	rof              uint8 // Rate Of Fire - How many frames between shots
	rofCounter       uint8
	mag              uint8   // Magazine - How many bullets the weapon can hold
	curMag           uint8   // How many bullets left
	bulletSize       float64 // How large of an area the weapon can hit
	image            *ebiten.Image
	fireImage        *ebiten.Image
	animationCounter uint8
}

func (w *Weapon) update() {
	if w.animationCounter > 0 {
		w.animationCounter--
	}
	if w.rofCounter > 0 {
		w.rofCounter--
	}
}

func (w *Weapon) draw(g *Game, screen *ebiten.Image, z float64) {
	var img *ebiten.Image
	if w.animationCounter == 0 {
		img = w.image
	} else {
		img = w.fireImage
	}

	ogW, ogH := img.Size()
	sW := screenWidth / 6.0 / float64(ogW)
	sH := screenHeight / 4.0 / float64(ogH)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)

	op.GeoM.Translate(screenWidth/2.0-(sW*float64(ogW))/2.0, screenHeight/8.0*7-sH*float64(ogH)+z)

	screen.DrawImage(img, &op)
}

func (w *Weapon) shoot(p *Player, enemies []*Enemy, bosses []*Boss, tiles []Tile) bool {
	var e *Enemy
	var b *Boss
	var dx, dy, dis, angle float64

	if w.rofCounter > 0 {
		return false
	}

	w.rofCounter = w.rof

	w.curMag--

	w.animationCounter = 5

	for i := len(bosses) - 1; i >= 0; i-- {
		b = bosses[i]

		// Check if behind a wall
		visible := true
		for j := 0; j < len(tiles); j++ {
			if tiles[j].check_line_in_tile(p.x, p.y, b.x, b.y) {
				visible = false
				break
			}
		}
		if !visible {
			continue
		}

		dx = b.x - p.x
		dy = b.y - p.y
		dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
		angle = to_degrees(math.Acos(dx / dis))

		if math.Asin(dy/dis) < 0 {
			angle = -angle
		}

		// How much to the left or right of the player the enemy is
		angle -= p.angle
		angle = bound_angle(angle)

		if angle < w.bulletSize || angle > 360-w.bulletSize {
			b.health -= w.damage
			return true
		}
	}

	for i := len(enemies) - 1; i >= 0; i-- {
		e = enemies[i]

		// Check if behind a wall
		visible := true
		for j := 0; j < len(tiles); j++ {
			if tiles[j].check_line_in_tile(p.x, p.y, e.x, e.y) {
				visible = false
				break
			}
		}
		if !visible {
			continue
		}

		dx = e.x - p.x
		dy = e.y - p.y
		dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
		angle = to_degrees(math.Acos(dx / dis))

		if math.Asin(dy/dis) < 0 {
			angle = -angle
		}

		// How much to the left or right of the player the enemy is
		angle -= p.angle
		angle = bound_angle(angle)

		if angle < w.bulletSize || angle > 360-w.bulletSize {
			e.health -= w.damage
			return true
		}
	}
	return false
}
