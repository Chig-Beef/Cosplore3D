package main

type Trigger struct {
	x      float64
	y      float64
	w      float64
	h      float64
	action Action
}

func (t *Trigger) check_collide(g *Game) {
	x := g.player.x
	y := g.player.y

	if t.x <= x && x <= t.x+t.w {
		if t.y <= y && y <= t.y+t.h {
			t.action(g)
		}
	}
}
