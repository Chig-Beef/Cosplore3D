package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Menu struct {
	buttons []*Button
}

func (m *Menu) update(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		for i := 0; i < len(m.buttons); i++ {
			if m.buttons[i].check_click(float64(g.curMousePos[0]), float64(g.curMousePos[1])) {
				m.buttons[i].onClick(g)
				break
			}
		}
	}
}

func (m *Menu) draw(screen *ebiten.Image, g *Game) {
	img := g.images["planet"]
	ogW, ogH := img.Size()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenWidth/3/ogW), float64(screenWidth/3/ogH))
	op.GeoM.Translate(float64(screenWidth/2-screenWidth/6), float64(screenHeight/2-screenWidth/6))
	op.ColorScale.ScaleAlpha(0.5)

	screen.DrawImage(img, &op)

	text.Draw(screen, "Cosplore3D", g.fonts["title"], screenWidth/2.0-240, screenHeight/2, color.RGBA{255, 255, 255, 255})

	for i := 0; i < len(m.buttons); i++ {
		m.buttons[i].draw(screen, g)
	}
}
