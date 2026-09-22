package monde

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/math32"
)

// Positions réservées sur la piste. L'équipe ennemie fait face au joueur.
// Elles ne déclenchent aucun combat : ce sont les repères pour l'intégration.
var PositionJoueurArene = math32.Vector3{X: 166, Y: -165, Z: .265}
var PositionSortieArene = math32.Vector3{X: 140, Y: -165, Z: .265}
var PositionsMonstresArene = []math32.Vector3{
	{X: 174, Y: -165, Z: .265},
	{X: 175, Y: -162, Z: .265},
	{X: 175, Y: -168, Z: .265},
}

// Appeler pendant le combat à la place de PlacerCameraPersonnage.
// La souris doit être libérée à l'entrée et son regard suspendu pendant le combat.
func PlacerCameraArene(vue *camera.Camera) {
	vue.SetPosition(170, -180, 9)
	vue.LookAt(&math32.Vector3{X: 170, Y: -165, Z: 1.5}, &math32.Vector3{Z: 1})
}

func EstDansArene(positionX float32, positionY float32) bool {
	centreX := float32(170)
	centreY := float32(-165)
	rayon := float32(25)

	differenceX := positionX - centreX
	differenceY := positionY - centreY

	distanceCarree := differenceX*differenceX + differenceY*differenceY

	rayonCarre := rayon * rayon

	return distanceCarree <= rayonCarre
}
