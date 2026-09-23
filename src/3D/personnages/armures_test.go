package personnages

import (
	"github.com/g3n/engine/math32"
	"testing"
)

// TestArmuresDetailleesEtEmplacements : Vérifie les pièces d'armure, leur rattachement et leur visibilité par emplacement.
func TestArmuresDetailleesEtEmplacements(t *testing.T) {
	p := Nouveau()
	if p.noeud.Rotation() != (math32.Vector3{}) || p.jambeDroite.Position().Z != .8 || p.brasDroit.Position().Z != 1.43 {
		t.Fatal("repères du personnage modifiés")
	}
	if len(p.chapeau.Children()) < 10 || len(p.tunique.Children()) < 25 || len(p.bottes) != 2 || len(p.protectionsBras) != 2 {
		t.Fatal("pièces détaillées manquantes")
	}
	for _, choix := range [][3]bool{{false, false, false}, {true, false, false}, {false, true, false}, {false, false, true}, {true, true, true}, {false, false, false}} {
		p.Equiper(choix[0], choix[1], choix[2])
		if p.chapeau.Visible() != choix[0] || p.tunique.Visible() != choix[1] {
			t.Fatal("équipement du corps incorrect")
		}
		for _, b := range p.bottes {
			if b.Visible() != choix[2] || len(b.Children()) < 19 {
				t.Fatal("bottes ou jambières incorrectes")
			}
		}
		for _, b := range p.protectionsBras {
			if b.Visible() != choix[1] || len(b.Children()) < 10 {
				t.Fatal("épaulières ou brassards incorrects")
			}
		}
	}
	// Les pièces restent attachées aux membres, pas à la racine immobile.
	if p.bottes[0].Parent() != p.jambeGauche || p.bottes[1].Parent() != p.jambeDroite || p.protectionsBras[0].Parent() != p.brasGauche || p.protectionsBras[1].Parent() != p.brasDroit {
		t.Fatal("armures détachées des articulations")
	}
}

// TestTenuesPNJConservees : Vérifie que les tenues du marchand et du forgeron restent présentes.
func TestTenuesPNJConservees(t *testing.T) {
	for _, nom := range []string{"marchand", "forgeron"} {
		p := NouveauAvecTenue(nom)
		if p.tunique.Visible() || p.bottes[0].Visible() || p.protectionsBras[0].Visible() {
			t.Fatal("armure du joueur visible sur un PNJ")
		}
		if p.chapeau.Visible() != (nom == "marchand") {
			t.Fatal("chapeau du PNJ modifié")
		}
	}
}
