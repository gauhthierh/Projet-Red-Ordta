package personnages

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

// Les formes de l'ancien modèle utilisent Y comme hauteur dans leur repère
// local. Seul le nœud visuel est tourné dans modele.go : le monde reste Z-up.
func boite(parent *core.Node, x, hauteur, avant, largeur, taille, profondeur float32, teinte *math32.Color) *graphic.Mesh {
	mesh := graphic.NewMesh(geometry.NewBox(largeur, taille, profondeur), matiereSurface(teinte))
	mesh.SetPosition(x, hauteur, avant)
	parent.Add(mesh)
	return mesh
}

func sphere(parent *core.Node, x, hauteur, avant, rayon float32, teinte *math32.Color) *graphic.Mesh {
	mesh := graphic.NewMesh(geometry.NewSphere(float64(rayon), 8, 6), matiereSurface(teinte))
	mesh.SetPosition(x, hauteur, avant)
	parent.Add(mesh)
	return mesh
}

func cylindre(parent *core.Node, x, hauteur, avant, rayon, taille float32, teinte *math32.Color) *graphic.Mesh {
	mesh := graphic.NewMesh(geometry.NewCylinder(float64(rayon), float64(taille), 10, 1, true, true), matiereSurface(teinte))
	mesh.SetPosition(x, hauteur, avant)
	parent.Add(mesh)
	return mesh
}
