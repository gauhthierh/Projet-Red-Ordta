// Package monde contient la scène 3D et sa boucle de jeu.
package monde

import (
	"math"
	"time"

	"ordta/3D/personnages"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

const (
	vitesseMarche    float32 = 4 // Unités parcourues par seconde.
	vitesseSprint    float32 = 6
	facteurDiagonale float32 = 0.7071
	rayonPersonnage  float32 = 0.3
)

// Lancer construit le monde, puis démarre sa boucle d'affichage.
func Lancer() {
	// Fenêtre G3N et racine de la scène 3D.
	ordta := app.App()
	scene := core.NewNode()

	// Gestion de la map
	CheminOBJ := "../assets/maps/red_world/red_world_map_3d.obj"
	CheminMTL := "../assets/maps/red_world/red_world_map_3d.mtl"

	donneesMonde, err := ChargerZones("../assets/maps/red_world/red_world_layout.json")
	if err != nil {
		panic(err)
	}

	spawn := donneesMonde.Spawn

	maps3d, err := ChargerMap3D(CheminOBJ, CheminMTL)

	if err != nil {
		panic(err)
	}
	scene.Add(maps3d)

	// Le package personnages fournit le modèle détaillé.
	// Son nœud de déplacement reste Z-up, comme le sol et la caméra.
	personnage3d := personnages.Nouveau()
	noeudPersonnage := personnage3d.Noeud()
	scene.Add(noeudPersonnage)
	noeudPersonnage.SetPosition(spawn[0], spawn[1], spawn[2])

	// Caméra suiveuse et lumière générale. G3N ouvre ici une fenêtre 800 x 600.
	vue := camera.New(800.0 / 600.0)
	scene.Add(vue)
	scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 1))
	haut := math32.NewVector3(0, 0, 1)

	// Boucle principale : mise à jour du jeu, puis affichage de chaque image.
	ordta.Gls().ClearColor(0.15, 0.25, 0.30, 1)
	ordta.Run(func(rendu *renderer.Renderer, tempsImage time.Duration) {
		positionActuelle := noeudPersonnage.Position()
		vitessePersonnage := vitesseMarche

		// Distance parcourue pendant cette frame.
		distance := vitessePersonnage * float32(tempsImage.Seconds())

		// Direction demandée par le joueur.
		directionX, directionY := float32(0), float32(0)

		if ordta.KeyState().Pressed(window.KeyW) {
			directionY += 1
		}

		if ordta.KeyState().Pressed(window.KeyS) {
			directionY -= 1
		}

		if ordta.KeyState().Pressed(window.KeyA) {
			directionX -= 1
		}

		if ordta.KeyState().Pressed(window.KeyD) {
			directionX += 1
		}
		if ordta.KeyState().Pressed(window.KeyLeftShift) {
			vitessePersonnage = vitesseSprint
			distance = vitessePersonnage * float32(tempsImage.Seconds())
		}

		// Deux touches simultanées ne doivent pas accélérer le personnage.
		if directionX != 0 && directionY != 0 {
			const diagonale = float32(0.70710678) // 1 / sqrt(2)
			directionX *= diagonale
			directionY *= diagonale
		}

		// Calcul de la nouvelle position dans le repère du monde.
		nouveauX := positionActuelle.X + directionX*distance
		nouveauY := positionActuelle.Y + directionY*distance

		// Application du déplacement.
		taillemap := donneesMonde.Taille
		obstacles := donneesMonde.Obstacles

		dansMap := estDansMap(nouveauX, nouveauY, taillemap, rayonPersonnage)
		toucheCercle := personnageToucheObstacleRond(nouveauX, nouveauY, rayonPersonnage, obstacles)
		toucheRectangle := personnageToucheRectangle(nouveauX, nouveauY, rayonPersonnage, obstacles)

		if dansMap && !toucheCercle && !toucheRectangle {
			noeudPersonnage.SetPosition(nouveauX, nouveauY, positionActuelle.Z)
		}

		// Garder la caméra derrière et au-dessus du personnage.
		positionCamera := noeudPersonnage.Position()
		vue.SetPosition(positionCamera.X, positionCamera.Y-10, positionCamera.Z+7)
		vue.LookAt(
			math32.NewVector3(positionCamera.X, positionCamera.Y, positionCamera.Z+0.9),
			haut,
		)

		// Animations du personnage.
		bouge := directionX != 0 || directionY != 0

		if bouge {
			angle := float32(math.Atan2(float64(directionX), float64(-directionY)))
			noeudPersonnage.SetRotationZ(angle)
		}

		personnage3d.Animer(float32(tempsImage.Seconds()), bouge)

		// Effacer l'image précédente, puis dessiner la nouvelle scène.
		ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
		rendu.Render(scene, vue)
	})
}
