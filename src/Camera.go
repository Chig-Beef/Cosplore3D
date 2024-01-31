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

	screen.DrawImage(output2D, &ebiten.DrawImageOptions{})
}

func (c *Camera) draw_world(level Level, screen *ebiten.Image) {
	// How many degrees between each ray
	angleDif := c.fov / screenWidth

	// Get all the tiles that can be hit (eg. code == 1)
	tiles := level.get_solid_tiles()

	newAngle := c.angle + c.fov/2.0
	bound_angle(&newAngle)

	for i := 0; i < screenWidth; i++ {
		tile, d := c.ray_cast(newAngle, tiles)
		newAngle -= angleDif
		bound_angle(&newAngle)
		if tile == (Tile{}) || d > screenHeight {
			continue
		}
		ebitenutil.DrawLine(screen, float64(i), d/2, float64(i), screenHeight-d/2, color.White)

		//fmt.Println("Cycle")
	}
}

func (c *Camera) ray_cast(angle float64, tiles []Tile) (Tile, float64) {
	x := c.x
	y := c.y

	var dx, dy float64

	// Move to the closest tile border
	if angle < 90 || angle > 270 {
		dy = float64(tileSize - (int(y) % tileSize))
	} else {
		dy = -float64(int(y) % tileSize)
	}
	if angle < 180 {
		dx = float64(tileSize - (int(x) % tileSize))
	} else {
		dx = -float64(int(x) % tileSize)
	}

	// Make sure we are going to a new tile border
	if dy == 0 {
		if angle < 90 || angle > 270 {
			dy = tileSize
		} else {
			dy = -tileSize
		}
	}
	if dx == 0 {
		if angle < 180 {
			dx = tileSize
		} else {
			dx = -tileSize
		}
	}

	if math.Abs(math.Cos(to_radians(angle))/dy) >= math.Abs(math.Sin(to_radians(angle))/dx) {
		dx = dy * math.Sin(to_radians(angle)) / math.Cos(to_radians(angle))
	} else {
		dy = dx * math.Cos(to_radians(angle)) / math.Sin(to_radians(angle))
	}

	y += dy
	x += dx
	d := math.Sqrt(math.Pow(x-c.x, 2) + math.Pow(y-c.y, 2))

	for d < c.dov {
		//fmt.Println(d)
		for i := 0; i < len(tiles); i++ {
			if tiles[i].check_hit(int(x), int(y)) {
				return tiles[i], d
			}
		}

		// Move to the closest tile border
		if angle < 90 || angle > 270 {
			dy = float64(tileSize - (int(y) % tileSize))
		} else {
			dy = -float64(int(y) % tileSize)
		}
		if angle < 180 {
			dx = float64(tileSize - (int(x) % tileSize))
		} else {
			dx = -float64(int(x) % tileSize)
		}
		//fmt.Println("a", dx, dy)

		if dy == 0 {
			if angle < 90 || angle > 270 {
				dy = tileSize
			} else {
				dy = -tileSize
			}
		}
		if dx == 0 {
			if angle < 180 {
				dx = tileSize
			} else {
				dx = -tileSize
			}
		}

		//fmt.Println("b", dx, dy)

		if math.Abs(math.Cos(to_radians(angle))/dy) >= math.Abs(math.Sin(to_radians(angle))/dx) {
			dx = dy * math.Sin(to_radians(angle)) / math.Cos(to_radians(angle))
		} else {
			dy = dx * math.Cos(to_radians(angle)) / math.Sin(to_radians(angle))
		}

		//fmt.Println("c", dx, dy)

		y += dy
		x += dx

		d = math.Sqrt(math.Pow(x-c.x, 2) + math.Pow(y-c.y, 2))
	}

	return Tile{}, 0
}

func (c *Camera) update_props(p *Player) {
	c.angle = p.angle
	c.x = p.x
	c.y = p.y
}
