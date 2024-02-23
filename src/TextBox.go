package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TextBox struct {
	text      string
	prefix    string
	x         float64
	y         float64
	w         float64
	h         float64
	bgColor   color.Color
	textcolor color.Color
	active    bool
}

func (tb *TextBox) draw(screen *ebiten.Image, g *Game) {
	ebitenutil.DrawRect(screen, tb.x, tb.y, tb.w, tb.h, tb.bgColor)

	text.Draw(screen, tb.prefix+tb.text, g.fonts["textBox"], int(tb.x+10), int(tb.y+tb.h/2), tb.textcolor)

	if tb.active {
		ebitenutil.DrawLine(screen, tb.x, tb.y+tb.h, tb.x+tb.w, tb.y+tb.h, color.White)
	}
}

func (tb *TextBox) check_click(g *Game) bool {
	x := g.curMousePos[0]
	y := g.curMousePos[1]

	if int(tb.x) <= x && x <= int(tb.x+tb.w) {
		if int(tb.y) <= y && y <= int(tb.y+tb.h) {
			tb.active = !tb.active
			return true
		}
	}

	return false
}

func (tb *TextBox) update() {
	if !tb.active {
		return
	}

	for i := 43; i < 53; i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key(i)) {
			tb.text += strconv.Itoa(i - 43)
			return
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(tb.text) > 0 {
			tb.text = tb.text[:len(tb.text)-1]
		}
	}
}
