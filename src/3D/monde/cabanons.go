package monde

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Les collisions générales des étals protègent l'intérieur réservé aux vendeurs.
func AjouterCabanons(scene *core.Node) {
	bois := material.NewStandard(&math32.Color{R: .8, G: .6, B: .4})
	tex, err := chargerTextureRepetee("../assets/maps/red_world/textures/wood_detailed.png")
	if err != nil {
		panic(err)
	}
	bois.AddTexture(tex)
	toit := material.NewStandard(&math32.Color{R: .65, G: .22, B: .16})
	piece := func(x, y, z, w, d, h float32, mat *material.Standard) {
		mesh := graphic.NewMesh(geometry.NewBox(w, d, h), mat)
		mesh.SetPosition(x, y, z)
		scene.Add(mesh)
	}
	for _, centre := range [][2]float32{{35, 14}, {48, 10}} {
		x, y := centre[0], centre[1]
		piece(x, y+2.3, 1.65, 7, .25, 3.3, bois)
		piece(x, y-2.5, .65, 7, .65, 1.1, bois)
		piece(x, y-2.5, 1.24, 7.4, .9, .12, bois)
		for _, dx := range []float32{-3.35, 3.35} {
			piece(x+dx, y-2.4, 1.7, .22, .22, 3.4, bois)
			piece(x+dx, y+2.3, 1.7, .22, .22, 3.4, bois)
		}
		piece(x, y, 3.55, 7.8, 5.9, .22, toit)
	}
}
