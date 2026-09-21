package monde

import "testing"

func TestTroncNonTraversablePendantUneImageLente(t *testing.T) {
	d := DonneesMonde{Taille: 600, Obstacles: []Collision{{Forme: "circle", Categorie: "tree", Rayon: .7}}}
	x, _ := deplacerAvecCollisions(-3, 0, 6, 0, .3, d)
	if x > -1 {
		t.Fatalf("tronc traversé : x=%v", x)
	}
}

func TestPontAutoriseSeulementLaRiviere(t *testing.T) {
	d := DonneesMonde{Taille: 600, Obstacles: []Collision{{Forme: "rectangle", Categorie: "river", Largeur: 50, Hauteur: 15}}, Ponts: []Collision{{Forme: "rectangle", Largeur: 25, Hauteur: 8, Rotation: 90}}}
	if !positionAutorisee(0, 0, .3, d) {
		t.Fatal("pont bloqué")
	}
	if positionAutorisee(8, 0, .3, d) {
		t.Fatal("eau traversable hors pont")
	}
	d.Obstacles = append(d.Obstacles, Collision{Forme: "circle", Categorie: "tree", Rayon: 1})
	if positionAutorisee(0, 0, .3, d) {
		t.Fatal("le pont annule une collision solide")
	}
}
