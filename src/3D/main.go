package main

import (
	"math"
	"time"

	"ordta/3D/personnages"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

const (
	// Taille d'une portion de terrain et rayon de la grille visible.
	taillePortion float32 = 100
	rayonGrille           = 1 // Une portion centrale et une de chaque côté : 3 x 3.
	coteGrille            = 2*rayonGrille + 1

	vitessePersonnage float32 = 3 // Unités parcourues par seconde.
	facteurDiagonale  float32 = 0.7071
)

func main() {
	// Fenêtre G3N et racine de la scène 3D.
	ordta := app.App()
	monde := core.NewNode()

	// Le package personnages fournit le modèle détaillé.
	// Son nœud de déplacement reste Z-up, comme le sol et la caméra.
	personnage3d := personnages.Nouveau()
	noeudPersonnage := personnage3d.Noeud()
	monde.Add(noeudPersonnage)

	// Terrain : neuf portions identiques, conservées dans une liste pour être
	// repositionnées autour du personnage pendant la partie.
	formeSol := geometry.NewPlane(taillePortion, taillePortion)
	matiereSol := material.NewStandard(math32.NewColor("DarkGreen"))
	sols := make([]*graphic.Mesh, 0, coteGrille*coteGrille)
	for x := -rayonGrille; x <= rayonGrille; x++ {
		for y := -rayonGrille; y <= rayonGrille; y++ {
			sol := graphic.NewMesh(formeSol, matiereSol)
			sol.SetPosition(float32(x)*taillePortion, float32(y)*taillePortion, 0)

			sols = append(sols, sol)
			monde.Add(sol)
		}
	}

	// Caméra suiveuse et lumière générale. G3N ouvre ici une fenêtre 800 x 600.
	vue := camera.New(800.0 / 600.0)
	monde.Add(vue)
	monde.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 1))
	haut := math32.NewVector3(0, 0, 1)

	// Boucle principale : mise à jour du jeu, puis affichage de chaque image.
	ordta.Gls().ClearColor(0.15, 0.25, 0.30, 1)
	ordta.Run(func(rendu *renderer.Renderer, tempsImage time.Duration) {
		// Lire ZQSD et construire une direction sur le plan du sol.
		directionX, directionY := float32(0), float32(0)
		if ordta.KeyState().Pressed(window.KeyW) {
			directionY++
		}
		if ordta.KeyState().Pressed(window.KeyS) {
			directionY--
		}
		if ordta.KeyState().Pressed(window.KeyA) {
			directionX--
		}
		if ordta.KeyState().Pressed(window.KeyD) {
			directionX++
		}

		// Deux touches simultanées ne doivent pas accélérer le personnage.
		if directionX != 0 && directionY != 0 {
			directionX *= facteurDiagonale
			directionY *= facteurDiagonale
		}
		distance := vitessePersonnage * float32(tempsImage.Seconds())
		noeudPersonnage.TranslateX(directionX * distance)
		noeudPersonnage.TranslateY(directionY * distance)

		// Recycler la grille autour de la portion occupée par le personnage.
		position := noeudPersonnage.Position()
		colonne := int(math.Round(float64(position.X / taillePortion)))
		ligne := int(math.Round(float64(position.Y / taillePortion)))
		for i, sol := range sols {
			decalageX := i/coteGrille - rayonGrille
			decalageY := i%coteGrille - rayonGrille
			sol.SetPosition(
				float32(colonne+decalageX)*taillePortion,
				float32(ligne+decalageY)*taillePortion,
				0,
			)
		}

		// Garder la caméra derrière et au-dessus du personnage.
		vue.SetPosition(position.X, position.Y-10, position.Z+7)
		vue.LookAt(
			math32.NewVector3(position.X, position.Y, position.Z+0.9),
			haut,
		)

		// Animations du personnage.
		bouge := directionX != 0 || directionY != 0
		if bouge {
			bouge = true
		} else {
			bouge = false
		}
		personnage3d.Animer(float32(tempsImage.Seconds()), bouge)

		// Effacer l'image précédente, puis dessiner la nouvelle scène.
		ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
		rendu.Render(monde, vue)
	})
}
