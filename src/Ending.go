package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ending struct {
	images  []*ebiten.Image
	counter uint16
	max     uint16
}

func (e *Ending) update() {
	e.counter++

	if e.counter == e.max {
		os.Exit(0)
	}
}

func (e *Ending) draw(screen *ebiten.Image) {
	img := e.images[(e.counter/7)%2]
	ogW, ogH := img.Size()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(screenWidth/float64(ogW), screenHeight/float64(ogH))
	op.ColorScale.ScaleAlpha(float32(e.max-e.counter) / float32(e.max))

	screen.DrawImage(img, &op)
}
