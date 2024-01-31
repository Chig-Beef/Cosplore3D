package main

type Tile struct {
	x    float64
	y    float64
	w    float64
	h    float64
	code uint8
}

func (t *Tile) check_hit(x, y int) bool {

	if x == int(t.x) || x == int(t.x+t.w) {
		if y >= int(t.y) && y <= int(t.y+t.h) {
			return true
		}
	}
	if y == int(t.y) || y == int(t.y+t.h) {
		if x >= int(t.x) && x <= int(t.x+t.w) {
			return true
		}
	}

	return false
}
