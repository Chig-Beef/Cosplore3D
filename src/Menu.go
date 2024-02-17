package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Menu struct {
	buttons []*Button
}

func (m *Menu) update(g *Game) {
	for i := 0; i < len(m.buttons); i++ {
		if m.buttons[i].check_click(float64(g.curMousePos[0]), float64(g.curMousePos[1])) {
			m.buttons[i].onClick(g, m)
			break
		}
	}
}

func (m *Menu) draw(screen *ebiten.Image, g *Game) {
	text.Draw(screen, "Cosplore3D", g.fonts["title"], screenWidth/2.0-240, 80, color.RGBA{255, 255, 255, 255})

	for i := 0; i < len(m.buttons); i++ {
		m.buttons[i].draw(screen, g)
	}
}
