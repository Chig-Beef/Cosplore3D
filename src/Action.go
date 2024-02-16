package main

type Action func(*Game)

func pickup_ammo(g *Game) {
	g.player.weapon.curMag += 10
	if g.player.weapon.curMag > g.player.weapon.mag {
		g.player.weapon.curMag = g.player.weapon.mag
	}
}

func pickup_cosmium(g *Game) {
	g.player.inventory = append(g.player.inventory, InvItem{"Cosmium", g.images["cosmium"]})
}

func rid_cosmium(g *Game) {
	for i := 0; i < len(g.player.inventory); i++ {
		if g.player.inventory[i].name == "Cosmium" {
			g.player.inventory = append(g.player.inventory[:i], g.player.inventory[i+1:]...)
			break
		}
	}
	l := g.levels[g.player.curLevel]
	for row := 0; row < len(l.data); row++ {
		for col := 0; col < len(l.data[row]); col++ {
			if l.data[row][col].code == 4 { // Cosplorer Reactor
				l.data[row][col].imageCols = g.imageColumns["cosplorerReactor"]
			}
		}
	}
}
