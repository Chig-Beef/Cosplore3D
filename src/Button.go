package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Button struct {
	x float64
	y float64
	w float64
	h float64

	bgColor   color.Color
	textColor color.Color

	text string

	onClick clickAction
}

type clickAction func(*Game)

func (b *Button) check_click(x, y float64) bool {
	return b.x <= x && x <= b.x+b.w && b.y <= y && y <= b.y+b.h
}

func (b *Button) draw(screen *ebiten.Image, g *Game) {
	ebitenutil.DrawRect(screen, b.x, b.y, b.w, b.h, b.bgColor)
	text.Draw(screen, b.text, g.fonts["btnText"], int(b.x+10), int(b.y+b.h/2+18), b.textColor)
}

func start_game(g *Game) {
	g.open_level("ankaran")
}

func open_settings(g *Game) {
	g.player.curLevel = "settings"
}

func open_menu(g *Game) {
	g.player.curLevel = "menu"
}
