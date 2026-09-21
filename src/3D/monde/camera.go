package monde

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

func ActiverRegardSouris(fenetre core.IDispatcher, angleHorizontal *float32, angleVertical *float32) {
	const sensibilite = float32(0.0025)

	var anciennePositionX float32
	var anciennePositionY float32
	premierMouvement := true

	fenetre.Subscribe(window.OnCursor, func(nomEvenement string, evenement interface{}) {
		mouvement := evenement.(*window.CursorEvent)

		// Le premier événement sert seulement à enregistrer
		// la position initiale de la souris.
		if premierMouvement {
			anciennePositionX = mouvement.Xpos
			anciennePositionY = mouvement.Ypos
			premierMouvement = false
			return
		}

		// Différence entre l'ancienne et la nouvelle position.
		deplacementX := mouvement.Xpos - anciennePositionX
		deplacementY := mouvement.Ypos - anciennePositionY

		anciennePositionX = mouvement.Xpos
		anciennePositionY = mouvement.Ypos

		// Souris vers la droite : rotation horizontale positive.
		*angleHorizontal += deplacementX * sensibilite

		// L'axe Y de la souris est inversé :
		// la valeur diminue lorsque la souris monte.
		*angleVertical -= deplacementY * sensibilite

		// Empêche la caméra de se retourner complètement.
		*angleVertical = math32.Clamp(*angleVertical, -1.45, 1.45)
	})
}

func PlacerCameraPersonnage(
	vue *camera.Camera,
	personnage *core.Node,
	angleHorizontal float32,
	angleVertical float32,
) {
	positionPersonnage := personnage.Position()

	cosVertical := math32.Cos(angleVertical)

	directionX := math32.Sin(angleHorizontal) * cosVertical
	directionY := math32.Cos(angleHorizontal) * cosVertical
	directionZ := math32.Sin(angleVertical)

	direction := math32.NewVector3(
		directionX,
		directionY,
		directionZ,
	)

	hauteurDesYeux := float32(1.70)

	positionDesYeux := math32.NewVector3(
		positionPersonnage.X,
		positionPersonnage.Y,
		positionPersonnage.Z+hauteurDesYeux,
	)

	decalageAvant := float32(0.20)

	positionCamera := math32.NewVector3(
		positionDesYeux.X+direction.X*decalageAvant,
		positionDesYeux.Y+direction.Y*decalageAvant,
		positionDesYeux.Z,
	)

	cible := math32.NewVector3(
		positionCamera.X+direction.X,
		positionCamera.Y+direction.Y,
		positionCamera.Z+direction.Z,
	)

	vue.SetPosition(
		positionCamera.X,
		positionCamera.Y,
		positionCamera.Z,
	)

	axeVertical := math32.NewVector3(0, 0, 1)

	vue.LookAt(cible, axeVertical)
}
