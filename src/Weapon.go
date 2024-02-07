package main

import "github.com/hajimehoshi/ebiten"

type Weapon struct {
	damage float64
	rof    uint8 // Rate Of Fire - How many frames between shots
	mag    uint8 // Magazine - How many bullets the weapon can hold
}

func (w *Weapon) draw(g *Game, screen *ebiten.Image) {
	img := g.images["gun"]
	ogW, ogH := img.Size()
	sW := screenWidth / 6.0 / float64(ogW)
	sH := screenHeight / 4.0 / float64(ogH)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)

	op.GeoM.Translate(screenWidth/2.0-(sW*float64(ogW))/2.0, screenHeight/8.0*7-sH*float64(ogH))

	screen.DrawImage(img, &op)
}
