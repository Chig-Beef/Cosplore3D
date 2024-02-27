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
	x              float64
	y              float64
	z              float64
	angle          float64
	camera         *Camera
	curLevel       string
	speed          float64
	haste          float64
	health         float64
	weapon         *Weapon
	weapons        []*Weapon
	curWeaponIndex int
	inventory      []InvItem
	bobCounter     uint8
}

func (p *Player) update(g *Game) {
	p.angle += float64(g.curMousePos[0]-g.prevMousePos[0]) * p.haste / 10.0 * g.sensitivity

	nx := p.x
	ny := p.y
	rx := p.x
	ry := p.y

	isMoving := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nx += math.Cos(to_radians(p.angle)) * p.speed
		ny += math.Sin(to_radians(p.angle)) * p.speed
		rx += math.Cos(to_radians(p.angle)) * p.speed * 2
		ry += math.Sin(to_radians(p.angle)) * p.speed * 2
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		nx -= math.Cos(to_radians(p.angle)) * p.speed
		ny -= math.Sin(to_radians(p.angle)) * p.speed
		rx -= math.Cos(to_radians(p.angle)) * p.speed * 2
		ry -= math.Sin(to_radians(p.angle)) * p.speed * 2
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		nx += math.Cos(to_radians(p.angle-90)) * p.speed
		ny += math.Sin(to_radians(p.angle-90)) * p.speed
		rx += math.Cos(to_radians(p.angle-90)) * p.speed * 2
		ry += math.Sin(to_radians(p.angle-90)) * p.speed * 2
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		nx -= math.Cos(to_radians(p.angle-90)) * p.speed
		ny -= math.Sin(to_radians(p.angle-90)) * p.speed
		rx -= math.Cos(to_radians(p.angle-90)) * p.speed * 2
		ry -= math.Sin(to_radians(p.angle-90)) * p.speed * 2
		isMoving = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		p.curWeaponIndex++
		if p.curWeaponIndex == len(p.weapons) {
			p.curWeaponIndex = 0
		}
		p.weapon = p.weapons[p.curWeaponIndex]
	}

	if isMoving {
		p.bobCounter += 5
		if p.bobCounter > 180 {
			p.bobCounter = 0
		}
		p.z = (math.Sin(to_radians(float64(p.bobCounter))) + 0.5) * 15
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
		if p.weapon.shoot(p, g.levels[p.curLevel].enemies, g.levels[p.curLevel].bosses, g.levels[p.curLevel].get_solid_tiles()) {
			g.play_effect("shoot")
		}
	}

	p.check_collide(g, g.levels[p.curLevel])

	p.weapon.update()
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
	p.weapon.draw(g, screen, p.z)

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

	text.Draw(screen, strconv.Itoa(int(p.weapon.curMag))+"/"+strconv.Itoa(int(p.weapon.mag))+": "+p.weapon.name, g.fonts["ammo"], screenWidth/3.2, screenHeight-10, color.RGBA{196, 32, 32, 255})

	for i := 0; i < len(p.inventory); i++ {
		p.inventory[i].draw(screen, float64(screenWidth/2+200+i*60), float64(screenHeight)/8.0*7+10)
	}
}
