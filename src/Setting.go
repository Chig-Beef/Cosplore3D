package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Setting struct {
	buttons   []*Button
	textBoxes map[string]*TextBox
}

func (s *Setting) update(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x := float64(g.curMousePos[0])
		y := float64(g.curMousePos[1])
		for _, b := range s.buttons {
			if b.check_click(x, y) {
				b.onClick(g)
			}
		}
		for _, tb := range s.textBoxes {
			tb.check_click(g)
		}
	}
	for _, tb := range s.textBoxes {
		tb.update()
	}
}

func (s *Setting) draw(screen *ebiten.Image, g *Game) {
	img := g.images["planet"]
	ogW, ogH := img.Size()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenWidth/3/ogW), float64(screenWidth/3/ogH))
	op.GeoM.Translate(float64(screenWidth/2-screenWidth/6), float64(screenHeight/2-screenWidth/6))
	op.ColorScale.ScaleAlpha(0.25)

	screen.DrawImage(img, &op)

	for _, b := range s.buttons {
		b.draw(screen, g)
	}
	for _, tb := range s.textBoxes {
		tb.draw(screen, g)
	}
}
