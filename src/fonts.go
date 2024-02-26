package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

func (g *Game) load_fonts() {
	g.fonts = make(map[string]font.Face)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	g.load_font("ammo", tt, dpi, screenHeight/8-20)
	g.load_font("title", tt, dpi, screenHeight/8)
	g.load_font("btnText", tt, dpi, screenHeight/12)
	g.load_font("textBox", tt, dpi, screenHeight/14)
}

func (g *Game) load_font(mName string, tt *sfnt.Font, dpi float64, size float64) {
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts[mName] = mplusNormalFont
}
