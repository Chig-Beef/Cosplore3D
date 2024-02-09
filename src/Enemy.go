package main

import (
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	x              float64
	y              float64
	w              float64
	h              float64
	angle          float64
	images         []*ebiten.Image
	target         *Player
	health         float64
	speed          float64
	damage         float64
	roa            uint8
	dov            float64
	attackRange    float64
	attackCooldown uint8
}

func (e *Enemy) update(g *Game) {
	if e.target != nil {
		e.follow_target()
		e.attack_target()
	} else {
		// If player close enough
		var dx, dy, dis float64

		dx = e.x - g.player.x
		dy = e.y - g.player.y
		dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

		if dis < e.dov {
			e.target = g.player
		}
	}
}

func (e *Enemy) attack_target() {
	var dx, dy, dis float64

	dx = e.x - e.target.x
	dy = e.y - e.target.y
	dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	if dis < e.attackRange {
		e.attackCooldown--
		if e.attackCooldown == 0 {
			e.target.health -= e.damage

			// Attempt to kill
			if e.target.health <= 0 {
				os.Exit(0)
			}

			e.attackCooldown = e.roa
		}
	}
}

func (e *Enemy) follow_target() {
	var dx, dy, dis, angle float64

	dx = e.x - e.target.x
	dy = e.y - e.target.y
	dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle = to_degrees(math.Acos(dx / dis))

	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	// How much to the left or right of the player the enemy is
	angle = bound_angle(angle)

	e.x -= math.Cos(to_radians(angle)) * e.speed
	e.y -= math.Sin(to_radians(angle)) * e.speed

	e.angle = angle
}

func (e *Enemy) draw(screen *ebiten.Image, c *Camera) {
	dx := e.x - c.x
	dy := e.y - c.y
	dis := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle := to_degrees(math.Acos(dx / dis))

	// In the behind 180 degrees
	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	angle -= c.angle

	angle = bound_angle(angle)

	if angle > c.fov/2.0 && angle < 360-c.fov/2.0 {
		return
	}

	lineX := (angle/(c.fov/2.0))*screenWidth/2.0 + screenWidth/2.0

	for lineX > screenWidth {
		lineX -= screenWidth
	}

	//ebitenutil.DrawLine(screen, lineX, 0, lineX, screenHeight, color.RGBA{255, 0, 0, 255})

	ogW, ogH := e.images[0].Size()
	sW := float64(ogW) / (dis / tileSize) * e.w / tileSize
	sH := float64(ogH) / (dis / tileSize) * e.h / tileSize

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)
	op.GeoM.Translate(lineX-(sW*float64(ogW)/2.0), screenHeight/2+sH*float64(ogH))
	screen.DrawImage(e.images[0], op)
}
