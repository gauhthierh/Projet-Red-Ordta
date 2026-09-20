package game3d

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"projet-red/internal/worldmap"
)

func (w *world) createVillage() {
	for _, id := range []string{"house_west", "house_north", "house_east"} {
		p := worldmap.LandmarkByID(id)
		house(w.scene, p.X, p.Z)
	}
	// Village cobbles and landmarks remain persistent; surrounding wilderness streams.
	for z := -6; z <= 32; z++ {
		for x := -1; x <= 1; x++ {
			box(w.scene, float32(x)*.70, .035, float32(z)*.72, .63, .065, .64, color(.43, .42, .36))
		}
	}
	for x := -31; x <= 31; x++ {
		for z := -1; z <= 1; z++ {
			box(w.scene, float32(x)*.72, .035, 4+float32(z)*.7, .65, .064, .64, color(.40, .41, .36))
		}
	}
	for _, p := range []struct{ x, z float32 }{{-15, 4}, {15, 4}, {0, 12}} {
		cylinder(w.scene, p.x, .78, p.z, .07, 1.55, wood)
		box(w.scene, p.x, 1.43, p.z, .9, .34, .11, wood)
		box(w.scene, p.x, 1.43, p.z+.065, .72, .17, .03, gold)
	}
	// Three authored destinations turn the roads into a small exploration loop.
	west := worldmap.LandmarkByID("wolf_camp")
	box(w.scene, west.X, 0.25, west.Z-1.4, 2.4, .5, 1.5, wood)
	for _, side := range []float32{-1, 1} {
		m := box(w.scene, west.X+side*.67, 1.05, west.Z-1.4, 1.55, .09, 1.9, roof)
		m.SetRotationZ(-side * .52)
	}
	crate(w.scene, west.X+1.4, .3, west.Z+1.1)
	cylinder(w.scene, west.X-1.3, .35, west.Z+1.1, .4, .25, stone)
	sphere(w.scene, west.X-1.3, .56, west.Z+1.1, .25, color(.95, .37, .08))
	east := worldmap.LandmarkByID("hill_cache")
	for _, p := range []struct{ x, z float32 }{{east.X - 1.1, east.Z - 1.3}, {east.X + 1.1, east.Z - 1.2}, {east.X, east.Z + 1.2}} {
		crate(w.scene, p.x, .3, p.z)
	}
	box(w.scene, east.X, .13, east.Z, 3.2, .22, 2.7, stone)
	south := worldmap.LandmarkByID("old_shrine")
	box(w.scene, south.X, .22, south.Z, 4.1, .42, 3.1, stone)
	for _, side := range []float32{-1, 1} {
		cylinder(w.scene, south.X+side*1.5, 1.3, south.Z, .3, 2.2, stone)
	}
	box(w.scene, south.X, 2.5, south.Z, 3.7, .36, .7, stone)
	sphere(w.scene, south.X, 2.93, south.Z, .36, gold)
	for _, x := range []float32{-6, -4} {
		cylinder(w.scene, x, 1.15, 0, .065, 2.3, wood)
	}
	for i := 0; i < 7; i++ {
		c := roof
		if i%2 == 0 {
			c = plaster
		}
		box(w.scene, -6.05+float32(i)*.35, 2.3, 0, .35, .1, 1.7, c).SetRotationX(.15)
	}
	box(w.scene, -5, .7, 0, 2.4, .18, 1.1, wood)
	for i := 0; i < 5; i++ {
		cylinder(w.scene, -5.7+float32(i)*.32, .9, 0, .10, .25, color(.45, .15+float32(i)*.08, .1))
		sphere(w.scene, -5.7+float32(i)*.32, 1.05, 0, .07, gold)
	}
	crate(w.scene, -6, .3, 1.2)
	crate(w.scene, -4.2, .3, 1.2)
	npc := newCharacter(false)
	npc.root.SetPosition(-5, 0, -1)
	npc.root.SetRotationY(.3)
	w.scene.Add(npc.root)
	box(w.scene, 5, .9, 0, 1.6, 1.8, 1.1, stone)
	box(w.scene, 5, .8, .565, 1.0, .95, .03, dark)
	for i := 0; i < 4; i++ {
		cone(w.scene, 4.65+float32(i)*.23, .65, .65, .15, .65, color(1, .35+float32(i)*.08, .04))
	}
	box(w.scene, 5, 2.6, 0, .6, 1.7, .65, stone)
	for i := 0; i < 5; i++ {
		box(w.scene, 5, 1.9+float32(i)*.32, 0, .68, .09, .73, color(.26, .28, .29))
	}
	cylinder(w.scene, 4.6, .3, 1.7, .35, .6, wood)
	box(w.scene, 4.6, .67, 1.7, 1.0, .20, .4, steel)
	box(w.scene, 5.7, .85, 1.5, .8, .16, .8, wood)
	glow := light.NewPoint(color(1, .45, .13), 2)
	glow.SetPosition(5, 1.2, .9)
	w.scene.Add(glow)
	smith := newCharacter(false)
	smith.root.SetPosition(6, 0, 1)
	smith.root.SetRotationY(-.6)
	w.scene.Add(smith.root)
	for _, x := range []float32{-2.2, 2.2} {
		cylinder(w.scene, x, 1.3, 3, .055, 2.6, wood)
		box(w.scene, x, 2.65, 3, .25, .38, .25, gold)
		lamp := light.NewPoint(color(1, .77, .36), .6)
		lamp.SetPosition(x, 2.7, 3)
		w.scene.Add(lamp)
	}
}

func house(parent *core.Node, x, z float32) {
	n := core.NewNode()
	n.SetPosition(x, 0, z)
	parent.Add(n)
	shadow(n, 0, 0, 3.3, 2.9)
	box(n, 0, 1.25, 0, 3, 2.5, 2.5, plaster)
	for _, side := range []float32{-1, 1} {
		m := box(n, side*.82, 2.90, 0, 1.95, .15, 3.0, roof)
		m.SetRotationZ(-side * .55)
	}
	box(n, 0, 3.38, 0, .16, .18, 3.1, wood)
	box(n, 1, 3.25, -.6, .4, 1.1, .4, stone)
	for _, xx := range []float32{-1.45, 0, 1.45} {
		box(n, xx, 1.25, 1.27, .13, 2.5, .10, wood)
	}
	for _, yy := range []float32{.12, 1.15, 2.43} {
		box(n, 0, yy, 1.27, 3, .12, .10, wood)
	}
	for _, xx := range []float32{-1, 1} {
		box(n, xx, 1.7, 1.30, .48, .55, .08, dark)
		box(n, xx, 1.7, 1.36, .035, .55, .04, gold)
		box(n, xx, 1.7, 1.36, .48, .035, .04, gold)
	}
	box(n, 0, .62, 1.33, .62, 1.15, .12, wood)
	sphere(n, .20, .6, 1.42, .04, gold)
	box(n, 0, .08, 1.6, .9, .16, .55, stone)
	// Timber lintel, flower pots and wall lantern add readable detail up close.
	box(n, 0, 1.29, 1.36, .87, .11, .17, wood)
	for _, side := range []float32{-1, 1} {
		box(n, side*1.1, .26, 1.62, .35, .29, .32, wood)
		for _, offset := range []float32{-.09, .09} {
			sphere(n, side*1.1+offset, .47, 1.62, .12, color(.20, .45, .16))
		}
	}
	box(n, 1.30, 1.20, 1.39, .17, .43, .19, wood)
	box(n, 1.30, 1.36, 1.42, .22, .20, .21, gold)
	for _, side := range []float32{-1, 1} {
		m := box(n, side*.72, 1.76, 1.34, .1, 1.45, .1, wood)
		m.SetRotationZ(side * .7)
	}
}
func crate(parent *core.Node, x, y, z float32) {
	box(parent, x, y, z, .6, .6, .6, wood)
	for _, yy := range []float32{y - .23, y + .23} {
		box(parent, x, yy, z+.31, .65, .08, .05, plaster)
	}
	box(parent, x, y, z+.32, .08, .6, .04, plaster)
}
