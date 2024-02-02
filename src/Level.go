package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
)

type Level struct {
	name string
	data [][]Tile
}

func load_levels(tileSize float64) map[string]Level {
	levels := make(map[string]Level)

	var levelData [][]string

	rawLevelData, err := os.ReadFile("maps/levels.json")
	if err != nil {
		return levels
	}

	err = json.Unmarshal(rawLevelData, &levelData)
	if err != nil {
		log.Fatal("failed to load `./maps/levels.json`, file may have been tampered with, reinstall advised")
	}

	// Get every level held in ./maps/
	for i := 0; i < len(levelData); i++ {
		fName := levelData[i][0]
		lName := levelData[i][1]
		rawLevel, err := os.ReadFile("maps/" + fName)
		if err != nil {
			return levels
		}

		// Now we have a list of strongs as the data
		rawRows := strings.Split(string(rawLevel), "\r\n")

		tiles := [][]Tile{}
		for row := 0; row < len(rawRows); row++ {
			tiles = append(tiles, []Tile{})
			for col := 0; col < len(rawRows[row]); col++ {
				code, err := strconv.Atoi(string(rawRows[row][col]))
				if err != nil {
					log.Fatal("failed to load a level correctly")
				}
				clr := get_color(uint8(code))
				tiles[row] = append(tiles[row], Tile{
					float64(col) * tileSize,
					float64(row) * tileSize,
					tileSize,
					tileSize,
					uint8(code),
					clr,
				})
			}
		}

		levels[lName] = Level{lName, tiles}
	}
	return levels
}

func get_color(code uint8) color.Color {
	switch code {
	case 0:
		return color.Black
	case 1:
		return color.RGBA{255, 255, 255, 255}
	case 2:
		return color.RGBA{255, 0, 0, 255}
	case 3:
		return color.RGBA{0, 255, 0, 255}
	case 4:
		return color.RGBA{0, 0, 255, 255}
	default:
		return color.Black
	}
}

func (l *Level) get_solid_tiles() []Tile {
	tiles := []Tile{}

	for row := 0; row < len(l.data); row++ {
		for col := 0; col < len(l.data); col++ {
			if l.data[row][col].code != 0 {
				tiles = append(tiles, l.data[row][col])
			}
		}
	}

	return tiles
}
