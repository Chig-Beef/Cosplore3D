package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Item struct {
	x      float64
	y      float64
	w      float64
	h      float64
	images []*ebiten.Image
	action Action
}

func (i *Item) check_collide(p *Player) bool {
	dx := i.x - p.x
	dy := i.y - p.y
	dis := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	return dis < float64(tileSize)/2
}

func (i *Item) draw(screen *ebiten.Image, c *Camera) {
	dx := i.x - c.x
	dy := i.y - c.y
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

	ogW, ogH := i.images[0].Size()
	sW := float64(ogW) / (dis / tileSize) * i.w / tileSize
	sH := float64(ogH) / (dis / tileSize) * i.h / tileSize

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)
	op.GeoM.Translate(lineX-(sW*float64(ogW)/2.0), screenHeight/2-sH*float64(ogH)/2)
	screen.DrawImage(i.images[0], op)
}
