package main

type Action func(*Game)

func pickup_ammo(g *Game) {
	g.player.weapon.curMag += 10
	if g.player.weapon.curMag > g.player.weapon.mag {
		g.player.weapon.curMag = g.player.weapon.mag
	}
}

func pickup_cosmium(g *Game) {
	g.player.inventory = append(g.player.inventory, "Cosmium")
}
