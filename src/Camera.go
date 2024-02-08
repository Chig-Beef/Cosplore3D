package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Camera struct {
	x     float64
	y     float64
	angle float64
	fov   float64 // Field Of View
	dov   float64 // Depth Of View
}

func (c *Camera) draw_2D(level Level, screen *ebiten.Image) {
	output2D := ebiten.NewImage(200, 200)
	output2D.Fill(color.RGBA{32, 32, 32, 255})

	var clr color.Color

	for row := 0; row < len(level.data); row++ {
		for col := 0; col < len(level.data[row]); col++ {
			if level.data[row][col].code == 0 {
				clr = color.Black
			} else {
				clr = color.White
			}

			ebitenutil.DrawRect(
				output2D,
				(level.data[row][col].x-c.x)/8+100,
				(level.data[row][col].y-c.y)/8+100,
				(level.data[row][col].w)/8,
				(level.data[row][col].h)/8,
				clr,
			)
		}
	}

	ebitenutil.DrawLine(output2D, 100+math.Cos(to_radians(c.angle))*10, 100+math.Sin(to_radians(c.angle))*10, 100, 100, color.RGBA{0, 0, 255, 255})

	tiles := level.get_solid_tiles()

	var dof int
	var rx, ry, ra, xo, yo, aTan, nTan, hx, hy, vx, vy, disT float64
	disV, disH := c.dov, c.dov

	ra = c.angle - c.fov/2.0
	ra = bound_angle(ra)

	for r := 0; r < int(c.fov); r++ {

		// Hor Line Check
		dof = 0
		aTan = -1 / math.Tan(to_radians(ra))

		if ra > 180 && ra != 360 { // Looking up
			ry = c.y - float64(int(c.y)%tileSize)
			rx = (c.y-ry)*aTan + c.x
			yo = -tileSize
			xo = -yo * aTan
		} else if ra < 180 && ra != 0 { // Looking down
			ry = c.y - float64(int(c.y)%tileSize) + tileSize
			rx = (c.y-ry)*aTan + c.x
			yo = tileSize
			xo = -yo * aTan
		} else {
			rx = c.x
			ry = c.y
			dof = 8
		}

		for dof < 8 {
			hit := false
			for i := 0; i < len(tiles); i++ {
				if tiles[i].check_hit(int(rx), int(ry)) {
					dof = 8
					hit = true
					hx = rx
					hy = ry
					disH = math.Sqrt(math.Pow(hx-c.x, 2) + math.Pow(hy-c.y, 2))
					break
				}
			}
			if hit {
				break
			}

			rx += xo
			ry += yo
			dof++
		}

		//ebitenutil.DrawLine(output2D, 100, 100, (hx-c.x)/8.0+100, (hy-c.y)/8.0+100, color.RGBA{255, 0, 0, 255})

		// Ver Line Check
		dof = 0
		nTan = -math.Tan(to_radians(ra))

		if ra > 90 && ra < 270 { // Looking left
			rx = c.x - float64(int(c.x)%tileSize)
			ry = (c.x-rx)*nTan + c.y
			xo = -tileSize
			yo = -xo * nTan
		} else if ra < 90 || ra > 270 { // Looking right
			rx = c.x - float64(int(c.x)%tileSize) + tileSize
			ry = (c.x-rx)*nTan + c.y
			xo = tileSize
			yo = -xo * nTan
		} else if ra == 90 || ra == 270 {
			rx = c.x
			ry = c.y
			dof = 8
		}

		for dof < 8 {
			hit := false
			for i := 0; i < len(tiles); i++ {
				if tiles[i].check_hit(int(rx), int(ry)) {
					dof = 8
					hit = true
					vx = rx
					vy = ry
					disV = math.Sqrt(math.Pow(vx-c.x, 2) + math.Pow(vy-c.y, 2))
					break
				}
			}
			if hit {
				break
			}

			rx += xo
			ry += yo
			dof++
		}

		//ebitenutil.DrawLine(output2D, 100, 100, (vx-c.x)/8.0+100, (vy-c.y)/8.0+100, color.RGBA{0, 255, 0, 255})

		if disV < disH {
			rx = vx
			ry = vy
			disT = disV
		} else {
			rx = hx
			ry = hy
			disT = disH
		}

		if disT <= c.dov && disT >= 0.5 {
			ebitenutil.DrawLine(output2D, 100, 100, (rx-c.x)/8.0+100, (ry-c.y)/8.0+100, color.RGBA{0, 255, 0, 255})
		}

		ra += 1
		ra = bound_angle(ra)
	}

	screen.DrawImage(output2D, &ebiten.DrawImageOptions{})
}

func (c *Camera) draw_world(level Level, screen *ebiten.Image) {

	tiles := level.get_solid_tiles()

	var ra float64

	ra = c.angle - c.fov/2.0
	ra = bound_angle(ra)

	ri := c.fov / float64(screenWidth)

	for r := 0; r < screenWidth; r++ {

		t, disT, rx, ry, isV := c.ray_cast(ra, tiles)

		if !t.init {
			ra += ri
			ra = bound_angle(ra)
			continue
		}

		a := c.angle - ra
		a = bound_angle(a)
		if a > 180 {
			a -= 360
		}

		disT *= math.Cos(to_radians(a))

		if disT > c.dov || disT < 0.5 {
			ra += ri
			ra = bound_angle(ra)
			continue
		}

		lineH := (float64(tileSize) * screenHeight) / disT
		if lineH == screenHeight {
			lineH = screenHeight
		}

		lineX := screenWidth - (a/c.fov*screenWidth + screenWidth/2.0)
		y := float64(screenHeight)/2.0 - lineH/2.0

		//ebitenutil.DrawLine(screen, lineX, float64(screenHeight)/2.0-lineH/2.0, lineX, float64(screenHeight)/2.0+lineH/2.0, get_color_with_distance(t.color, disT))

		imgs := *t.imageCols

		var subX int

		if isV {
			subX = int(ry) % tileSize

			if ra < 90 || ra > 270 {
				subX = tileSize - 1 - subX
			}
		} else {
			subX = int(rx) % tileSize

			if ra > 180 {
				subX = tileSize - 1 - subX
			}
		}

		tx := int(float64(subX) * float64(len(imgs)) / float64(tileSize))

		img := imgs[tx]
		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(1, lineH/float64(img.Bounds().Dy()))
		op.GeoM.Translate(lineX, y)
		//op.ColorM.ChangeHSV(0, 255, float64(fastInvSqrt(float32(disT)/float32(tileSize))))
		//op.ColorM.ChangeHSV(0, 0, 1/math.Sqrt(disT/float64(tileSize)))
		op.ColorM.ChangeHSV(0, 1, float64(fastInvSqrt(float32(int(disT)/tileSize+1))))

		screen.DrawImage(img, &op)

		ra += ri
		ra = bound_angle(ra)
	}
}

func (c *Camera) ray_cast(ra float64, tiles []Tile) (Tile, float64, float64, float64, bool) {
	var dof int
	var maxDof int = int(c.dov) / tileSize
	var rx, ry, xo, yo, aTan, nTan, hx, hy, vx, vy float64
	disV, disH := c.dov, c.dov
	var disT float64
	var vt, ht, tt Tile
	var hitV, hitH bool

	// Hor Line Check
	dof = 0
	aTan = -1 / math.Tan(to_radians(ra))

	if ra > 180 && ra != 360 { // Looking up
		ry = c.y - float64(int(c.y)%tileSize)
		rx = (c.y-ry)*aTan + c.x
		yo = -tileSize
		xo = -yo * aTan
	} else if ra < 180 && ra != 0 { // Looking down
		ry = c.y - float64(int(c.y)%tileSize) + tileSize
		rx = (c.y-ry)*aTan + c.x
		yo = tileSize
		xo = -yo * aTan
	} else {
		rx = c.x
		ry = c.y
		dof = maxDof
	}

	for dof < maxDof {
		hitH = false
		for i := 0; i < len(tiles); i++ {
			if tiles[i].check_hit(int(rx), int(ry)) {
				dof = 8
				hitH = true
				hx = rx
				hy = ry
				disH = math.Sqrt(math.Pow(hx-c.x, 2) + math.Pow(hy-c.y, 2))
				ht = tiles[i]
				break
			}
		}
		if hitH {
			break
		}

		rx += xo
		ry += yo
		dof++
	}

	// Ver Line Check
	dof = 0
	nTan = -math.Tan(to_radians(ra))

	if ra > 90 && ra < 270 { // Looking left
		rx = c.x - float64(int(c.x)%tileSize)
		ry = (c.x-rx)*nTan + c.y
		xo = -tileSize
		yo = -xo * nTan
	} else if ra < 90 || ra > 270 { // Looking right
		rx = c.x - float64(int(c.x)%tileSize) + tileSize
		ry = (c.x-rx)*nTan + c.y
		xo = tileSize
		yo = -xo * nTan
	} else if ra == 90 || ra == 270 {
		rx = c.x
		ry = c.y
		dof = maxDof
	}

	for dof < maxDof {
		hitV = false
		for i := 0; i < len(tiles); i++ {
			if tiles[i].check_hit(int(rx), int(ry)) {
				dof = 8
				hitV = true
				vx = rx
				vy = ry
				disV = math.Sqrt(math.Pow(vx-c.x, 2) + math.Pow(vy-c.y, 2))
				vt = tiles[i]
				break
			}
		}
		if hitV {
			break
		}

		rx += xo
		ry += yo
		dof++
	}

	var isV bool

	if disV < disH && hitV {
		rx = vx
		ry = vy
		disT = disV
		tt = vt
		isV = true
	} else {
		rx = hx
		ry = hy
		disT = disH
		tt = ht
		isV = false
	}

	return tt, disT, rx, ry, isV

}

func (c *Camera) update_props(p *Player) {
	c.angle = p.angle
	c.x = p.x
	c.y = p.y
}
