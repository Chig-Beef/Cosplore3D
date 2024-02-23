package main

import (
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
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
}

func create_new_enemy(g *Game, col, row int, w, h float64, images []*ebiten.Image, health, speed, damage float64, roa uint8, dov, attackRange float64) *Enemy {
	return &Enemy{
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
	}
}

func (e *Enemy) set_difficulty(d float64) {
	e.health *= d
	e.speed *= d
	e.damage *= d
	e.roa = uint8(float64(e.roa) / d)
	e.dov *= d
	e.attackRange *= d
}

func (e *Enemy) update(g *Game, tiles []Tile) {
	if e.target != nil {
		e.follow_target(tiles)
		e.attack_target(g)
	} else {
		// If player close enough
		if e.get_distance(g.player.x, g.player.y, true) < e.dov {
			e.target = g.player
		}
	}
}

func (e *Enemy) attack_target(g *Game) {
	if e.get_distance(e.target.x, e.target.y, true) < e.attackRange {
		e.attackCooldown--
		if e.attackCooldown == 0 {
			e.target.health -= e.damage

			g.play_effect("playerHurt")

			// Attempt to kill
			if e.target.health <= 0 {
				os.Exit(0)
			}

			e.attackCooldown = e.roa
		}
	}
}

func (e *Enemy) get_distance(x, y float64, useSqrt bool) float64 {
	dx := e.x - x
	dy := e.y - y
	if useSqrt { // If this ever gets slow, just use fastinvsqrt, and see if it approximates well enough
		return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	} else { // Eliminates a SQRT call, improving performance
		return math.Pow(dx, 2) + math.Pow(dy, 2)
	}
}

func (e *Enemy) follow_target(tiles []Tile) {
	var dx, dy, dis, angle float64

	dx = e.x - e.target.x
	dy = e.y - e.target.y
	dis = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle = to_degrees(math.Acos(dx / dis))

	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	// How much to the left or right of the player the enemy is
	angle = bound_angle(angle)

	nx := e.x
	ny := e.y
	rx := e.x
	ry := e.y

	nx -= math.Cos(to_radians(angle)) * e.speed
	ny -= math.Sin(to_radians(angle)) * e.speed
	rx -= math.Cos(to_radians(angle)) * e.speed * 2
	ry -= math.Sin(to_radians(angle)) * e.speed * 2
	hit := false
	for i := 0; i < len(tiles); i++ {
		if tiles[i].check_point_in_tile(rx, ry) {
			hit = true
			break
		}
	}
	if !hit {
		e.x = nx
		e.y = ny
	}

	e.angle = angle
}

func (e *Enemy) draw(screen *ebiten.Image, c *Camera, tiles []Tile) {
	for i := 0; i < len(tiles); i++ {
		if tiles[i].check_line_in_tile(e.x, e.y, c.x, c.y) {
			// Can't see through a wall
			return
		}
	}

	dx := e.x - c.x
	dy := e.y - c.y
	dis := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	angle := to_degrees(math.Acos(dx / dis))

	// In the behind 180 degrees
	if math.Asin(dy/dis) < 0 {
		angle = -angle
	}

	vangle := angle - e.angle

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
		img = e.images[0]
	} else if vangle > 45 && vangle <= 135 {
		img = e.images[1]
	} else if vangle > 135 && vangle <= 225 {
		img = e.images[2]
	} else {
		img = e.images[3]
	}

	ogW, ogH := img.Size()
	sW := ((float64(tileSize) * screenHeight) * (e.w / tileSize) / dis) / float64(ogW)
	sH := ((float64(tileSize) * screenHeight) * (e.h / tileSize) / dis) / float64(ogH)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sW, sH)
	op.GeoM.Translate(lineX-(sW*float64(ogW)/2.0), screenHeight/2+((float64(tileSize)*screenHeight)/dis/2)-(sH*float64(ogH)))
	screen.DrawImage(img, op)
}

func (g *Game) kull_enemies() {
	var e *Enemy
	aliveEnemies := []int{}
	for i := 0; i < len(g.levels[g.player.curLevel].enemies); i++ {
		e = g.levels[g.player.curLevel].enemies[i]
		if e.health > 0 {
			aliveEnemies = append(aliveEnemies, i)
		} else {
			g.play_effect("enemyDeath")
		}
	}
	newEnemies := make([]*Enemy, len(aliveEnemies))
	for i := 0; i < len(aliveEnemies); i++ {
		newEnemies[i] = g.levels[g.player.curLevel].enemies[aliveEnemies[i]]
	}
	g.levels[g.player.curLevel].enemies = newEnemies
}
