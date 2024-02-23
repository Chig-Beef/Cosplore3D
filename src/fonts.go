package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func (g *Game) load_fonts() {
	g.fonts = make(map[string]font.Face)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight/8 - 20,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["ammo"] = mplusNormalFont

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight / 8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["title"] = mplusNormalFont

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight / 12,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["btnText"] = mplusNormalFont

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    screenHeight / 14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.fonts["textBox"] = mplusNormalFont
}
