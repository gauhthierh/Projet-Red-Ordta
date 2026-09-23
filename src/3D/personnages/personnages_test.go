package personnages

import (
	"math"
	"testing"

	"github.com/g3n/engine/math32"
)

// TestAncienPersonnageZDirect : Vérifie le repère Z vertical et les pivots de l'ancien modèle sans conversion globale.
func TestAncienPersonnageZDirect(t *testing.T) {
	p := NouveauAncienModele()
	if p.Noeud() == nil || p.jambeGauche == nil || p.jambeDroite == nil || p.brasGauche == nil || p.brasDroit == nil {
		t.Fatal("le personnage ou ses articulations sont absents")
	}
	if p.noeud.Rotation() != (math32.Vector3{}) {
		t.Fatal("la racine ne doit plus convertir les axes")
	}
	if p.jambeDroite.Position().Z != .8 || p.brasDroit.Position().Z != 1.42 {
		t.Fatal("les pivots doivent être placés directement sur l'axe Z")
	}
	if len(p.noeud.Children()) != 25 || len(p.jambeDroite.Children()) != 3 || len(p.brasDroit.Children()) != 7 {
		t.Fatal("des pièces du personnage ont disparu")
	}

	// Un sommet du personnage doit rester au-dessus du sol sans rotation du modèle.
	modele := lireModele()
	for _, groupe := range [][]piece{modele.Corps, modele.JambeGauche, modele.JambeDroite, modele.BrasGauche, modele.BrasDroit} {
		for _, partie := range groupe {
			if len(partie.Sommets)%3 != 0 || len(partie.Normales) != len(partie.Sommets) || len(partie.UV)/2 != len(partie.Sommets)/3 {
				t.Fatal("coordonnées du modèle statique incohérentes")
			}
			for _, index := range partie.Indices {
				if int(index) >= len(partie.Sommets)/3 {
					t.Fatal("indice de sommet hors du modèle statique")
				}
			}
		}
	}
	maxZ := float32(0)
	for _, piece := range modele.Corps {
		for i := 2; i < len(piece.Sommets); i += 3 {
			maxZ = float32(math.Max(float64(maxZ), float64(piece.Sommets[i])))
		}
	}
	if maxZ < 2 {
		t.Fatalf("hauteur du personnage trop faible : %v", maxZ)
	}
}
