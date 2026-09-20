package game3d

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"math"
)

type character struct {
	root       *core.Node
	legs, arms [2]*core.Node
	blade      *graphic.Mesh
	phase      float32
}

func newCharacter(goblin bool) *character {
	a := &character{root: core.NewNode()}
	bodyColor, faceColor := cloth, skin
	if goblin {
		bodyColor = wood
		faceColor = color(.30, .53, .17)
	}
	// Torso, neck, breastplate, belt, buckle and layered shoulder armour.
	box(a.root, 0, 1.10, 0, .62, .73, .34, bodyColor)
	box(a.root, 0, 1.25, .19, .46, .40, .08, steel)
	cylinder(a.root, 0, 1.55, 0, .12, .19, faceColor)
	box(a.root, 0, .83, 0, .68, .12, .40, wood)
	box(a.root, 0, .83, .23, .13, .13, .07, gold)
	for i, x := range []float32{-.24, .24} {
		a.legs[i] = core.NewNode()
		a.legs[i].SetPosition(x, .8, 0)
		a.root.Add(a.legs[i])
		box(a.legs[i], 0, -.24, 0, .22, .48, .24, bodyColor)
		box(a.legs[i], 0, -.60, .07, .26, .28, .40, wood)
		box(a.legs[i], 0, -.49, .20, .12, .07, .04, gold)
		a.arms[i] = core.NewNode()
		a.arms[i].SetPosition(x*1.65, 1.42, 0)
		a.root.Add(a.arms[i])
		sphere(a.arms[i], 0, 0, 0, .22, steel)
		box(a.arms[i], 0, -.24, 0, .19, .42, .21, bodyColor)
		box(a.arms[i], 0, -.46, 0, .19, .15, .23, wood)
		sphere(a.arms[i], 0, -.58, .04, .115, faceColor)
	}
	head := sphere(a.root, 0, 1.82, 0, .29, faceColor)
	head.SetScale(1, 1.12, .9)
	// Eyes, brows, nose, mouth and ears are modelled, not painted placeholders.
	for _, x := range []float32{-.105, .105} {
		box(a.root, x, 1.86, .243, .08, .065, .028, color(.94, .89, .70))
		box(a.root, x, 1.86, .265, .035, .045, .018, dark)
		box(a.root, x, 1.925, .24, .12, .035, .03, wood)
		sphere(a.root, x*2.65, 1.82, 0, .07, faceColor)
	}
	box(a.root, 0, 1.80, .265, .075, .10, .07, faceColor)
	box(a.root, 0, 1.70, .24, .10, .025, .025, wood)
	if goblin {
		for _, x := range []float32{-.35, .35} {
			ear := cone(a.root, x, 1.90, 0, .12, .40, faceColor)
			ear.SetRotationZ(-x * 2)
		}
		for _, x := range []float32{-.07, .07} {
			cone(a.root, x, 1.70, .27, .035, .13, color(.92, .87, .66))
		}
	} else {
		sphere(a.root, 0, 2.02, -.02, .28, wood).SetScale(1, .55, 1)
		for _, x := range []float32{-.22, .22} {
			box(a.root, x, 1.91, -.03, .09, .27, .30, wood)
		}
		box(a.root, 0, 1.12, -.24, .58, .82, .08, color(.48, .12, .10)) // cape
		box(a.root, .13, 1.36, -.31, .36, .42, .16, wood)               // backpack
	}
	// Sword and hilt are attached to the right hand; shield to the left.
	a.blade = box(a.arms[1], 0, -.3, .18, .075, .72, .045, steel)
	box(a.arms[1], 0, -.64, .18, .31, .055, .08, gold)
	box(a.arms[1], 0, -.74, .18, .075, .18, .07, wood)
	shield := sphere(a.arms[0], -.1, -.35, .16, .29, wood)
	shield.SetScale(1, 1.3, .22)
	sphere(a.arms[0], -.1, -.35, .23, .075, steel)
	if goblin {
		a.root.SetScale(.9, .9, .9)
	}
	return a
}

func (a *character) animate(dt float32, moving bool, attack float32) {
	stride := float32(0)
	if moving {
		a.phase += dt * 9
		stride = float32(math.Sin(float64(a.phase))) * .55
	}
	a.legs[0].SetRotationX(stride)
	a.legs[1].SetRotationX(-stride)
	a.arms[0].SetRotationX(-stride * .6)
	a.arms[1].SetRotationX(stride * .6)
	if attack > 0 {
		a.arms[1].SetRotationX(-float32(math.Sin(float64((1-attack/.28)*math.Pi))) * 1.6)
	}
}
