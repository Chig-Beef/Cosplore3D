package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Camera struct {
	x     float64
	y     float64
	angle float64
	fov   float64
	dov   float64
}

func (c *Camera) draw_2D(level Level, screen *ebiten.Image) {
	output2D, _ := ebiten.NewImage(200, 200, ebiten.FilterDefault)
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

	//ray_length := 64.0

	//fmt.Println(to_degrees(math.Atan((-float64(screenWidth)/2.0 + float64(0)) / (float64(screenWidth) / 2.0))))
	//fmt.Println(to_degrees(math.Atan((-float64(screenWidth)/2.0 + float64(screenWidth)) / (float64(screenWidth) / 2.0))))

	/*for i := 0; i < screenWidth; i += 30 {
		newAngle := c.angle - to_degrees(math.Atan((float64(i)-float64(screenWidth)/2.0)/(float64(screenWidth)/2.0)))
		bound_angle(&newAngle)

		ebitenutil.DrawLine(output2D, 100+math.Sin(to_radians(newAngle))*ray_length, 100+math.Cos(to_radians(newAngle))*ray_length, 100, 100, color.White)
	}*/

	tiles := level.get_solid_tiles()

	var dof int
	var rx, ry, ra, xo, yo, aTan, nTan, hx, hy, vx, vy, disT float64
	disV, disH := 10_000.0, 10_000.0

	ra = c.angle - c.fov/2.0
	bound_angle(&ra)

	for r := 0; r < int(c.fov); r++ {

		// Hor Line Check
		dof = 0
		aTan = -1 / math.Tan(to_radians(ra))

		if ra > 180 { // Looking up
			ry = c.y - float64(int(c.y)%tileSize)
			rx = (c.y-ry)*aTan + c.x
			yo = -tileSize
			xo = -yo * aTan
		} else if ra < 180 { // Looking down
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
		} else if ra == 180 || ra == 0 {
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

		if disT != 10_000.0 {
			ebitenutil.DrawLine(output2D, 100, 100, (rx-c.x)/8.0+100, (ry-c.y)/8.0+100, color.RGBA{0, 255, 0, 255})
		}

		ra += 1
		bound_angle(&ra)
	}

	screen.DrawImage(output2D, &ebiten.DrawImageOptions{})
}

func (c *Camera) draw_world(level Level, screen *ebiten.Image) {

	tiles := level.get_solid_tiles()

	var ra float64

	ra = c.angle - c.fov/2.0
	bound_angle(&ra)

	ri := c.fov / float64(screenWidth)

	for r := 0; r < int(screenWidth); r++ {

		t, disT := c.ray_cast(ra, tiles)

		a := c.angle - ra
		bound_angle(&a)
		if a > 180 {
			a -= 360
		}

		disT *= math.Cos(to_radians(a))

		if disT > 5_000 || disT < 0.5 {
			continue
		}

		lineH := (float64(tileSize) * screenHeight) / disT
		if lineH == screenHeight {
			lineH = screenHeight
		}

		lineX := screenWidth - (a/c.fov*screenWidth + screenWidth/2.0)

		ebitenutil.DrawLine(screen, lineX, float64(screenHeight)/2.0-lineH/2.0, lineX, float64(screenHeight)/2.0+lineH/2.0, get_color_with_distance(t.color, disT))

		ra += ri
		bound_angle(&ra)
	}
}

func (c *Camera) ray_cast(ra float64, tiles []Tile) (Tile, float64) {
	var dof int
	var rx, ry, xo, yo, aTan, nTan, hx, hy, vx, vy float64
	disV, disH := 10_000.0, 10_000.0
	var disT float64
	var vt, ht, tt Tile

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
				ht = tiles[i]
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
				vt = tiles[i]
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

	if disV < disH {
		//rx = vx
		//ry = vy
		disT = disV
		tt = vt
	} else {
		//rx = hx
		//ry = vy
		disT = disH
		tt = ht
	}

	return tt, disT

}

func (c *Camera) update_props(p *Player) {
	c.angle = p.angle
	c.x = p.x
	c.y = p.y
}
