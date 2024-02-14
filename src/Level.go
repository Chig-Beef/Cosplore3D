package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	name           string
	data           [][]Tile
	enemies        []*Enemy
	items          []*Item
	floorColor     color.RGBA
	skyColor       color.RGBA
	playerStartPos [2]float64
	progressors    []Progressor
	fullyLoaded    bool
}

func (l Level) draw_floor_sky(screen *ebiten.Image) {
	screen.Fill(l.skyColor)
	ebitenutil.DrawRect(screen, 0, screenHeight/2, screenWidth, screenHeight/2, l.floorColor)
}

func load_levels(g *Game, tileSize float64) map[string]*Level {
	levels := make(map[string]*Level)

	var levelData [][]string

	rawLevelData, err := os.ReadFile("assets/maps/levels.json")
	if err != nil {
		return levels
	}

	err = json.Unmarshal(rawLevelData, &levelData)
	if err != nil {
		log.Fatal("failed to load `./assets/maps/levels.json`, file may have been tampered with, reinstall advised")
	}

	var px, py float64

	// Get every level held in ./maps/
	for i := 0; i < len(levelData); i++ {
		fName := levelData[i][0]
		lName := levelData[i][1]
		rawLevel, err := os.ReadFile("assets/maps/" + fName)
		if err != nil {
			return levels
		}

		// Now we have a list of strongs as the data
		rawRows := strings.Split(string(rawLevel), "\r\n")

		floorColor, skyColor := get_fs_color(rawRows[0])

		rawRows = rawRows[1:]

		tiles := [][]Tile{}
		enemies := []*Enemy{}
		items := []*Item{}
		progressors := []Progressor{}

		for row := 0; row < len(rawRows); row++ {
			tiles = append(tiles, []Tile{})
			rawRow := strings.Split(rawRows[row], ",")
			for col := 0; col < len(rawRow); col++ {
				code, err := strconv.Atoi(string(rawRow[col]))
				if err != nil {
					log.Fatal("failed to load a level correctly")
				}

				// Player
				if code == 3 {
					code = 0
					px = tileSize * (float64(col) + 0.5)
					py = tileSize * (float64(row) + 0.5)
				}
				// Progressor
				if code == 4 {
					code = 0
					progressors = append(progressors, Progressor{
						tileSize * (float64(col)),
						tileSize * (float64(row)),
						tileSize,
						tileSize,
						has_cosmium,
						"cosplorer",
					})
				}
				// Enemy
				if code == 9 { // Blob
					code = 0
					enemies = append(enemies, &Enemy{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						100,
						70,
						0,
						[]*ebiten.Image{
							g.images["blobFront"],
							g.images["blobLeft"],
							g.images["blobBack"],
							g.images["blobRight"],
						},
						nil,
						100,
						1,
						1,
						60,
						5 * tileSize,
						10,
						60,
					})
				} else if code == 8 { // Infected Crewmate
					code = 0
					enemies = append(enemies, &Enemy{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						100,
						180,
						0,
						[]*ebiten.Image{
							g.images["crewmateFront"],
							g.images["crewmateLeft"],
							g.images["crewmateBack"],
							g.images["crewmateRight"],
						},
						nil,
						120,
						2,
						1.5,
						30,
						5 * tileSize,
						10,
						30,
					})
				}

				// Item
				if code == 7 {
					code = 0
					items = append(items, &Item{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						50,
						50,
						[]*ebiten.Image{g.images["cosmium"]},
						pickup_cosmium,
					})
				} else if code == 6 {
					code = 0
					items = append(items, &Item{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						50,
						50,
						[]*ebiten.Image{g.images["ammo"]},
						pickup_ammo,
					})
				}

				tiles[row] = append(tiles[row], Tile{
					true,
					float64(col) * tileSize,
					float64(row) * tileSize,
					tileSize,
					tileSize,
					uint8(code),
					g.imageColumns[get_tile_image(uint8(code))],
				})
			}
		}

		levels[lName] = &Level{lName, tiles, enemies, items, floorColor, skyColor, [2]float64{px, py}, progressors, false}
	}

	return levels
}

func get_tile_image(code uint8) string {
	switch code {
	case 0:
		return ""
	case 1:
		return "cosplorerWall"
	case 2:
		return "ankaranWall"
	case 5:
		return "cosplorerComputer"
	default:
		return "ankaranWall"
	}
}

func (l *Level) get_solid_tiles() []Tile {
	tiles := []Tile{}

	for row := 0; row < len(l.data); row++ {
		for col := 0; col < len(l.data[row]); col++ {
			if l.data[row][col].code != 0 {
				tiles = append(tiles, l.data[row][col])
			}
		}
	}

	return tiles
}

func get_fs_color(data string) (color.RGBA, color.RGBA) {
	splitData := strings.Split(data, "|")

	if len(splitData) != 2 {
		log.Fatal("need a floor and sky color")
	}

	colorArray := [3]uint8{0, 0, 0}

	fRaw := strings.Split(splitData[0], ",")
	if len(fRaw) != 3 {
		log.Fatal("3 values in color needed, RGB")
	}
	for i := 0; i < 3; i++ {
		c, err := strconv.Atoi(fRaw[i])
		if err != nil {
			log.Fatal("invalid number for color")
		}
		colorArray[i] = uint8(c)
	}
	floorColor := color.RGBA{colorArray[0], colorArray[1], colorArray[2], 255}

	sRaw := strings.Split(splitData[1], ",")
	if len(fRaw) != 3 {
		log.Fatal("3 values in color needed, RGB")
	}
	for i := 0; i < 3; i++ {
		c, err := strconv.Atoi(sRaw[i])
		if err != nil {
			log.Fatal("invalid number for color")
		}
		colorArray[i] = uint8(c)
	}
	skyColor := color.RGBA{colorArray[0], colorArray[1], colorArray[2], 255}

	return floorColor, skyColor
}

func (l *Level) draw_items_and_enemies(screen *ebiten.Image, c *Camera) {
	tiles := l.get_solid_tiles()
	for i := 0; i < len(l.enemies); i++ {
		l.enemies[i].draw(screen, c, tiles)
	}
	for i := 0; i < len(l.items); i++ {
		l.items[i].draw(screen, c, tiles)
	}
}

func (g *Game) open_level(levelName string) {
	g.player.curLevel = levelName

	if !g.levels[levelName].fullyLoaded {
		apply_image_colums_to_tiles(g, g.levels[levelName])
	}

	g.player.x = g.levels[levelName].playerStartPos[0]
	g.player.y = g.levels[levelName].playerStartPos[1]
}
