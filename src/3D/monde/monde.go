// Package monde contient la scène 3D et sa boucle de jeu.
package monde

import (
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
	vitesseMarche   float32 = 6 // Unités parcourues par seconde.
	vitesseSprint   float32 = 8
	rayonPersonnage float32 = 0.3
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

	// Caméra
	cameraSimulation := camera.New(1920.0 / 1080.0)

	ConfigurerFenetreSimulation(
		ordta,
		cameraSimulation,
	)

	VerrouillerSourisSimulation()
	angleHorizontal := float32(0)
	angleVertical := float32(0)

	ActiverRegardSouris(ordta, &angleHorizontal, &angleVertical)

	// Lumière dans le jeu
	scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 1))

	// Boucle principale : mise à jour du jeu, puis affichage de chaque image.
	ordta.Gls().ClearColor(0.15, 0.25, 0.30, 1)
	ordta.Run(func(rendu *renderer.Renderer, tempsImage time.Duration) {
		positionActuelle := noeudPersonnage.Position()
		vitessePersonnage := vitesseMarche

		// Distance parcourue pendant cette frame.
		distance := vitessePersonnage * float32(tempsImage.Seconds())

		// Direction demandée par le joueur.
		commandeAvant, commandeDroite := float32(0), float32(0)

		if ordta.KeyState().Pressed(window.KeyW) {
			commandeAvant += 1
		}

		if ordta.KeyState().Pressed(window.KeyS) {
			commandeAvant -= 1
		}

		if ordta.KeyState().Pressed(window.KeyA) {
			commandeDroite -= 1
		}

		if ordta.KeyState().Pressed(window.KeyD) {
			commandeDroite += 1
		}
		if ordta.KeyState().Pressed(window.KeyLeftShift) {
			vitessePersonnage = vitesseSprint
			distance = vitessePersonnage * float32(tempsImage.Seconds())
		}

		directionX, directionY := CalculerDirectionPersonnage(angleHorizontal, commandeAvant, commandeDroite)

		// Calcul de la nouvelle position dans le repère du monde.
		nouveauX := positionActuelle.X + directionX*distance
		nouveauY := positionActuelle.Y + directionY*distance

		// Les deux axes sont vérifiés séparément : le personnage glisse ainsi
		// naturellement le long d'un mur au lieu de rester bloqué en diagonale.
		positionFinaleX := positionActuelle.X
		positionFinaleY := positionActuelle.Y

		if positionAutorisee(nouveauX, positionActuelle.Y, rayonPersonnage, donneesMonde) {
			positionFinaleX = nouveauX
		}
		if positionAutorisee(positionFinaleX, nouveauY, rayonPersonnage, donneesMonde) {
			positionFinaleY = nouveauY
		}

		noeudPersonnage.SetPosition(positionFinaleX, positionFinaleY, positionActuelle.Z)
		deplacementEffectue := positionFinaleX != positionActuelle.X || positionFinaleY != positionActuelle.Y

		// Caméra
		PlacerCameraPersonnage(cameraSimulation, noeudPersonnage, angleHorizontal, angleVertical)

		// Animations du personnage.
		bouge := deplacementEffectue
		angle := math32.Pi - angleHorizontal
		noeudPersonnage.SetRotationZ(angle)
		personnage3d.Animer(float32(tempsImage.Seconds()), bouge)

		// Effacer l'image précédente, puis dessiner la nouvelle scène.
		ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
		rendu.Render(scene, cameraSimulation)
	})
}
