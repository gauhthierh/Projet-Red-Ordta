package monde

import (
	"fmt"
	"ordta/3D/monstres"
	"ordta/library"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

func ChargerMontres(scene *core.Node, concurrents []library.EnnemiCombat) ([]*monstres.Monstre, error) {
	modeles := []*monstres.Monstre{}

	if len(concurrents) > len(PositionsMonstresArene) {
		return nil, fmt.Errorf(
			"pas assez de positions pour %d concurrents",
			len(concurrents),
		)
	}

	for index, concurrent := range concurrents {
		chemin := "../assets/models/monstres/" +
			concurrent.Modele3D

		modele, err := monstres.Charger(chemin)
		if err != nil {
			return nil, fmt.Errorf(
				"charger %s : %w",
				concurrent.Modele3D,
				err,
			)
		}

		position := PositionsMonstresArene[index]
		noeud := modele.Noeud()

		noeud.SetPosition(
			position.X,
			position.Y,
			position.Z,
		)

		noeud.SetRotationZ(-math32.Pi / 2)

		scene.Add(noeud)
		modeles = append(modeles, modele)
	}

	return modeles, nil
}
