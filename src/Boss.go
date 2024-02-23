package main

import (
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Boss struct {
	x              float64
	y              float64
	w              float64
	h              float64
	angle          float64
	images         []*ebiten.Image
	target         *Player
	health         float64
	speed          float64
	damage         float64
	roa            uint8
	dov            float64
	attackRange    float64
	attackCooldown uint8
	power          Power
	rop            uint8
	powerCooldown  uint8
	visible        bool
}

type Power func(*Game, *Boss)

func create_new_boss(g *Game, col, row int, w, h float64, images []*ebiten.Image, health, speed, damage float64, roa uint8, dov, attackRange float64, power Power, rop uint8) *Boss {
	return &Boss{
		tileSize * (float64(col) + 0.5),
		tileSize * (float64(row) + 0.5),
		w,
		h,
		0,
		images,
		nil,
		health,
		speed,
		damage,
		roa,
		dov * tileSize,
		attackRange,
		roa,
		power,
		rop,
		rop,
		true,
	}
}

func (b *Boss) set_difficulty(d float64) {
	b.health *= d
	b.speed *= d
	b.damage *= d
	b.roa = uint8(float64(b.roa) / d)
	b.rop = uint8(float64(b.rop) / d)
	b.dov *= d
	b.attackRange *= d
}

func (b *Boss) update(g *Game, tiles []Tile) {
	if b.target != nil {
		b.follow_target(tiles)
		b.attack_target()
		b.powerCooldown--
		if b.powerCooldown <= 0 {
			b.powerCooldown = b.rop
			if b.power != nil {
				b.power(g, b)
			}
		}
	} else {
		// If player close enough
		if b.get_distance(g.player.x, g.player.y, true) < b.dov {
			b.target = g.player
		}
	}
}

func (b *Boss) attack_target() {
	if b.get_distance(b.target.x, b.target.y, true) < b.attackRange {
		b.attackCooldown--
		if b.attackCooldown == 0 {
			b.target.health -= b.damage

			// Attempt to kill
			if b.target.health <= 0 {
				os.Exit(0)
			}

			b.attackCooldown = b.roa
		}
	}
}

func (b *Boss) get_distance(x, y float64, useSqrt bool) float64 {
	dx := b.x - x
	dy := b.y - y
	if useSqrt { // If this ever gets slow, just use fastinvsqrt, and see if it approximates well enough
		return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	} else { // Eliminates a SQRT call, improving performance
		return math.Pow(dx, 2) + math.Pow(dy, 2)
	}
}

func (b *Boss) follow_target(tiles []Tile) {
	var dx, dy, dis, angle float64

	dx = b.x - b.target.x
	dy = b.y - b.target.y
	dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle = to_degrees(math.Acos(dx / dis))

	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	// How much to the left or right of the player the enemy is
	angle = bound_angle(angle)

	nx := b.x
	ny := b.y
	rx := b.x
	ry := b.y

	nx -= math.Cos(to_radians(angle)) * b.speed
	ny -= math.Sin(to_radians(angle)) * b.speed
	rx -= math.Cos(to_radians(angle)) * b.speed * 2
	ry -= math.Sin(to_radians(angle)) * b.speed * 2
	hit := false
	for i := 0; i < len(tiles); i++ {
		if tiles[i].check_point_in_tile(rx, ry) {
			hit = true
			break
		}
	}
	if !hit {
		b.x = nx
		b.y = ny
	}

	b.angle = angle
}

func (b *Boss) draw(screen *ebiten.Image, c *Camera, tiles []Tile) {
	if !b.visible {
		return
	}

	for i := 0; i < len(tiles); i++ {
		if tiles[i].check_line_in_tile(b.x, b.y, c.x, c.y) {
			// Can't see through a wall
			return
		}
	}

	dx := b.x - c.x
	dy := b.y - c.y
	dis := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle := to_degrees(math.Acos(dx / dis))

	// In the behind 180 degrees
	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	vangle := angle - b.angle

	angle -= c.angle

	angle = bound_angle(angle)

	dis *= math.Cos(to_radians(angle))

	if angle > c.fov/2.0 && angle < 360-c.fov/2.0 {
		return
	}

	lineX := (angle/(c.fov/2.0))*screenWidth/2.0 + screenWidth/2.0

	for lineX > screenWidth {
		lineX -= screenWidth
	}

	//ebitenutil.DrawLine(screen, lineX, 0, lineX, screenHeight, color.RGBA{255, 0, 0, 255})

	vangle = bound_angle(vangle)
	var img *ebiten.Image
	if vangle <= 45 || vangle >= 315 {
		img = b.images[0]
	} else if vangle > 45 && vangle <= 135 {
		img = b.images[1]
	} else if vangle > 135 && vangle <= 225 {
		img = b.images[2]
	} else {
		img = b.images[3]
	}

	ogW, ogH := img.Size()
	sW := ((float64(tileSize) * screenHeight) * (b.w / tileSize) / dis) / float64(ogW)
	sH := ((float64(tileSize) * screenHeight) * (b.h / tileSize) / dis) / float64(ogH)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)
	op.GeoM.Translate(lineX-(sW*float64(ogW)/2.0), screenHeight/2+((float64(tileSize)*screenHeight)/dis/2)-(sH*float64(ogH)))
	screen.DrawImage(img, op)
}

func (g *Game) kull_bosses() {
	var b *Boss
	aliveBosses := []int{}
	for i := 0; i < len(g.levels[g.player.curLevel].bosses); i++ {
		b = g.levels[g.player.curLevel].bosses[i]
		if b.health > 0 {
			aliveBosses = append(aliveBosses, i)
		}
	}
	newBosses := make([]*Boss, len(aliveBosses))
	for i := 0; i < len(aliveBosses); i++ {
		newBosses[i] = g.levels[g.player.curLevel].bosses[aliveBosses[i]]
	}
	g.levels[g.player.curLevel].bosses = newBosses
}

func change_visibility(g *Game, b *Boss) {
	b.visible = !b.visible
}

func spawn_crawler(g *Game, b *Boss) {
	g.levels[g.player.curLevel].enemies = append(g.levels[g.player.curLevel].enemies, create_new_enemy(
		g,
		int(b.x/tileSize),
		int(b.y/tileSize),
		100,
		70,
		[]*ebiten.Image{
			g.images["crawlerFront"],
			g.images["crawlerLeft"],
			g.images["crawlerBack"],
			g.images["crawlerRight"],
		},
		100,
		4,
		4,
		18,
		7,
		8,
	),
	)
}
