package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) load_images() {
	g.images = make(map[string]*ebiten.Image)

	load_image(g, "blob1", "blobFront")
	load_image(g, "blob2", "blobRight")
	load_image(g, "blob3", "blobBack")
	load_image(g, "blob4", "blobLeft")

	load_image(g, "heart", "heart")
	load_image(g, "gun", "gun")
	load_image(g, "cosplorerWall", "cosplorerWall")
	load_image(g, "ankaranWall", "ankaranWall")
	load_image(g, "cosmium", "cosmium")
	load_image(g, "ammo", "ammo")
}

func load_image(g *Game, fName string, mName string) {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/" + fName + ".png")
	if err != nil {
		log.Fatal(err)
	}
	g.images[mName] = img
}
