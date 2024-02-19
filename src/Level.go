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
	bosses         []*Boss
	items          []*Item
	floorColor     color.RGBA
	skyColor       color.RGBA
	playerStartPos [2]float64
	progressors    []Progressor
	triggers       []Trigger
	audio          string
	fullyLoaded    bool
}

func (l *Level) update_progressors_and_triggers(g *Game) {
	for i := 0; i < len(l.progressors); i++ {
		l.progressors[i].check_collide(g)
	}
	for i := 0; i < len(l.triggers); i++ {
		l.triggers[i].check_collide(g)
	}
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
		audio := rawRows[1]

		rawRows = rawRows[2:]

		tiles := [][]Tile{}
		enemies := []*Enemy{}
		bosses := []*Boss{}
		items := []*Item{}
		progressors := []Progressor{}
		triggers := []Trigger{}

		for row := 0; row < len(rawRows); row++ {
			tiles = append(tiles, []Tile{})
			rawRow := strings.Split(rawRows[row], ",")
			for col := 0; col < len(rawRow); col++ {
				code, err := strconv.Atoi(string(rawRow[col]))
				if err != nil {
					log.Fatal("failed to load a level correctly")
				}

				switch code {
				case 10: // Player
					code = 0
					px = tileSize * (float64(col) + 0.5)
					py = tileSize * (float64(row) + 0.5)
				case 20: // Progressor (Has Cosmium)
					code = 0
					progressors = append(progressors, Progressor{
						tileSize * (float64(col)),
						tileSize * (float64(row)),
						tileSize,
						tileSize,
						has_cosmium,
						"cosplorer",
					})
				case 21: // Progressor (Has No Cosmium)
					code = 0
					progressors = append(progressors, Progressor{
						tileSize * (float64(col)),
						tileSize * (float64(row)),
						tileSize,
						tileSize,
						has_no_cosmium,
						"enikoko",
					})
				case 22: // Progressor (Has Cosmium)
					code = 0
					progressors = append(progressors, Progressor{
						tileSize * (float64(col)),
						tileSize * (float64(row)),
						tileSize,
						tileSize,
						has_cosmium,
						"schmeltool",
					})
				case 30: // Ammo
					code = 0
					items = append(items, &Item{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						50,
						50,
						[]*ebiten.Image{g.images["ammo"]},
						pickup_ammo,
					})
				case 31: // Cosmium
					code = 0
					items = append(items, &Item{
						tileSize * (float64(col) + 0.5),
						tileSize * (float64(row) + 0.5),
						50,
						50,
						[]*ebiten.Image{g.images["cosmium"]},
						pickup_cosmium,
					})
				case 40: // Blob
					code = 0
					enemies = append(enemies, create_new_enemy(
						col,
						row,
						100,
						70,
						[]*ebiten.Image{
							g.images["blobFront"],
							g.images["blobLeft"],
							g.images["blobBack"],
							g.images["blobRight"],
						},
						100,
						1,
						1,
						60,
						5,
						10,
					),
					)
				case 41: // Infected Crewmate
					code = 0
					enemies = append(enemies, create_new_enemy(
						col,
						row,
						100,
						180,
						[]*ebiten.Image{
							g.images["crewmateFront"],
							g.images["crewmateLeft"],
							g.images["crewmateBack"],
							g.images["crewmateRight"],
						},
						120,
						2,
						1.5,
						30,
						5,
						20,
					),
					)
				case 42: // Trash Crawler
					code = 0
					enemies = append(enemies, create_new_enemy(
						col,
						row,
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
				case 43: // Enikoko Beast
					code = 0
					bosses = append(bosses, create_new_boss(
						col,
						row,
						200,
						280,
						[]*ebiten.Image{
							g.images["beastFront"],
							g.images["beastLeft"],
							g.images["beastBack"],
							g.images["beastRight"],
						},
						500,
						5,
						3,
						24,
						5,
						10,
						change_visibility,
						60,
					),
					)
				case 44: // Crawler Champion
					code = 0
					bosses = append(bosses, create_new_boss(
						col,
						row,
						250,
						200,
						[]*ebiten.Image{
							g.images["championFront"],
							g.images["championLeft"],
							g.images["championBack"],
							g.images["championRight"],
						},
						500,
						7,
						5,
						24,
						5,
						10,
						spawn_crawler,
						255,
					),
					)
				case 50: // Trigger (rid Cosmium)
					code = 0
					triggers = append(triggers, Trigger{
						tileSize * (float64(col)),
						tileSize * (float64(row)),
						tileSize,
						tileSize,
						rid_cosmium,
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

		levels[lName] = &Level{lName, tiles, enemies, bosses, items, floorColor, skyColor, [2]float64{px, py}, progressors, triggers, audio, false}
	}

	return levels
}

func get_tile_image(code uint8) string {
	switch code {
	case 0:
		return ""
	case 1:
		return "ankaranWall"
	case 2:
		return "cosplorerWall"
	case 3:
		return "cosplorerComputer"
	case 4:
		return "cosplorerReactorEmpty"
	case 5:
		return "enikokoWall"
	case 6:
		return "schmeltoolWall"
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

func (l *Level) draw_items_enemies_and_bosses(screen *ebiten.Image, c *Camera) {
	tiles := l.get_solid_tiles()
	l.sort_bosses_by_distance(c)
	l.sort_enemies_by_distance(c)
	l.sort_items_by_distance(c)
	for i := 0; i < len(l.enemies); i++ {
		l.enemies[i].draw(screen, c, tiles)
	}
	for i := 0; i < len(l.items); i++ {
		l.items[i].draw(screen, c, tiles)
	}
	for i := 0; i < len(l.bosses); i++ {
		l.bosses[i].draw(screen, c, tiles)
	}
}

func (l *Level) sort_enemies_by_distance(c *Camera) {
	for i := 0; i < len(l.enemies); i++ {
		for j := 0; j < len(l.enemies)-i-1; j++ {
			if l.enemies[j].get_distance(c.x, c.y, false) < l.enemies[j+1].get_distance(c.x, c.y, false) {
				temp := l.enemies[j]
				l.enemies[j] = l.enemies[j+1]
				l.enemies[j+1] = temp
			}
		}
	}
}

func (l *Level) sort_bosses_by_distance(c *Camera) {
	for i := 0; i < len(l.bosses); i++ {
		for j := 0; j < len(l.bosses)-i-1; j++ {
			if l.bosses[j].get_distance(c.x, c.y, false) < l.bosses[j+1].get_distance(c.x, c.y, false) {
				temp := l.bosses[j]
				l.bosses[j] = l.bosses[j+1]
				l.bosses[j+1] = temp
			}
		}
	}
}

func (l *Level) sort_items_by_distance(c *Camera) {
	for i := 0; i < len(l.items); i++ {
		for j := 0; j < len(l.items)-i-1; j++ {
			if l.items[j].get_distance(c.x, c.y, false) < l.items[j+1].get_distance(c.x, c.y, false) {
				temp := l.items[j]
				l.items[j] = l.items[j+1]
				l.items[j+1] = temp
			}
		}
	}
}

func (g *Game) open_level(levelName string) {
	g.player.curLevel = levelName

	if !g.hasLoadedImageColumns {
		g.create_image_columns([]string{
			"cosplorerWall",
			"ankaranWall",
			"cosplorerComputer",
			"cosplorerReactorEmpty",
			"cosplorerReactor",
			"enikokoWall",
			"schmeltoolWall",
		})
		g.apply_image_columns_to_tiles(g.levels[g.player.curLevel])
	}

	if !g.levels[levelName].fullyLoaded {
		g.apply_image_columns_to_tiles(g.levels[levelName])
	}

	g.player.x = g.levels[levelName].playerStartPos[0]
	g.player.y = g.levels[levelName].playerStartPos[1]

	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	g.play_audio(g.levels[levelName].audio)
}
