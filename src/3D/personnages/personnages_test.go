package personnages

import (
	"math"
	"testing"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

func TestPersonnageEtOrientation(t *testing.T) {
	p := Nouveau()
	if p.Noeud() == nil || p.jambeGauche == nil || p.jambeDroite == nil || p.brasGauche == nil || p.brasDroit == nil {
		t.Fatal("le personnage n'est pas construit complètement")
	}
	if ecart := math.Abs(float64(p.modele.Rotation().X - math32.Pi/2)); ecart > 0.0001 {
		t.Fatalf("orientation du modèle : écart de %f radian", ecart)
	}
	// Le Y local de l'ancien modèle devient bien le Z vertical de la scène.
	repere := core.NewNode()
	repere.SetPosition(0, 1, 0)
	p.modele.Add(repere)
	p.noeud.UpdateMatrixWorld()
	var position math32.Vector3
	repere.WorldPosition(&position)
	if math.Abs(float64(position.Y)) > 0.0001 || math.Abs(float64(position.Z-1)) > 0.0001 {
		t.Fatalf("axe vertical incorrect après adaptation : %+v", position)
	}
}
