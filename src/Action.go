package main

type Action func(*Game)

func giveAmmo(g *Game) {
	g.player.weapon.curMag += 10
	if g.player.weapon.curMag > g.player.weapon.mag {
		g.player.weapon.curMag = g.player.weapon.mag
	}
}
