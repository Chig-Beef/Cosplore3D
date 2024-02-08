package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	x      float64
	y      float64
	w      float64
	h      float64
	images []*ebiten.Image
	target Player
	health float64
	speed  float64
	damgae float64
	roa    uint8
}

func (e *Enemy) update() {

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
