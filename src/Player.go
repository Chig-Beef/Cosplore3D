package main

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Player struct {
	x         float64
	y         float64
	angle     float64
	camera    *Camera
	curLevel  string
	speed     float64
	haste     float64
	health    float64
	weapon    Weapon
	inventory []string
}

func (p *Player) update(g *Game) {
	p.angle += float64(g.curMousePos[0]-g.prevMousePos[0]) * p.haste / 10.0

	nx := p.x
	ny := p.y
	rx := p.x
	ry := p.y

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nx += math.Cos(to_radians(p.angle)) * p.speed
		ny += math.Sin(to_radians(p.angle)) * p.speed
		rx += math.Cos(to_radians(p.angle)) * p.speed * 2
		ry += math.Sin(to_radians(p.angle)) * p.speed * 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nx -= math.Cos(to_radians(p.angle)) * p.speed
		ny -= math.Sin(to_radians(p.angle)) * p.speed
		rx -= math.Cos(to_radians(p.angle)) * p.speed * 2
		ry -= math.Sin(to_radians(p.angle)) * p.speed * 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		nx += math.Cos(to_radians(p.angle-90)) * p.speed
		ny += math.Sin(to_radians(p.angle-90)) * p.speed
		rx += math.Cos(to_radians(p.angle-90)) * p.speed * 2
		ry += math.Sin(to_radians(p.angle-90)) * p.speed * 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		nx -= math.Cos(to_radians(p.angle-90)) * p.speed
		ny -= math.Sin(to_radians(p.angle-90)) * p.speed
		rx -= math.Cos(to_radians(p.angle-90)) * p.speed * 2
		ry -= math.Sin(to_radians(p.angle-90)) * p.speed * 2
	}

	hit := false
	tiles := g.levels[p.curLevel].get_solid_tiles()
	for i := 0; i < len(tiles); i++ {
		if tiles[i].check_point_in_tile(rx, ry) {
			hit = true
			break
		}
	}
	if !hit {
		p.x = nx
		p.y = ny
	}

	p.angle = bound_angle(p.angle)

	p.camera.update_props(p)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && p.weapon.curMag > 0 {
		p.weapon.shoot(p, g.levels[p.curLevel].enemies, g.levels[p.curLevel].get_solid_tiles())
	}

	p.check_collide(g, g.levels[p.curLevel])
}

func (p *Player) check_collide(g *Game, l *Level) {
	aliveItems := []int{}
	for i := 0; i < len(l.items); i++ {
		if l.items[i].check_collide(p) {
			if l.items[i].action != nil {
				l.items[i].action(g)
			}
		} else {
			aliveItems = append(aliveItems, i)
		}
	}
	newItems := make([]*Item, len(aliveItems))

	for i := 0; i < len(aliveItems); i++ {
		newItems[i] = l.items[aliveItems[i]]
	}

	l.items = newItems
}

func (p *Player) draw(g *Game, screen *ebiten.Image) {
	p.camera.draw_world(*g.levels[p.curLevel], screen)

	//p.camera.draw_2D(g.levels[p.curLevel], screen)
}

func (p *Player) draw_hud(g *Game, screen *ebiten.Image) {
	p.weapon.draw(g, screen)

	ebitenutil.DrawRect(screen, 0, float64(screenHeight)/8.0*7, screenWidth, float64(screenHeight)/8.0, color.Gray{128})

	heartImg := g.images["heart"]
	ogW, ogH := heartImg.Size()
	newW := (float64(screenHeight)/8.0 - 20)
	newH := (float64(screenHeight)/8.0 - 20)

	for i := 0; i < 5; i++ {

		if p.health < float64((i+1)*20) {
			continue
		}
		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(newW/float64(ogW), newH/float64(ogH))

		op.GeoM.Translate(10*float64(i+1)+newW*float64(i), float64(screenHeight)/8.0*7+10)

		screen.DrawImage(heartImg, &op)
	}

	text.Draw(screen, strconv.Itoa(int(p.weapon.curMag))+"/"+strconv.Itoa(int(p.weapon.mag)), g.fonts["ammo"], screenWidth/2.0, screenHeight-10, color.RGBA{196, 32, 32, 255})
}
