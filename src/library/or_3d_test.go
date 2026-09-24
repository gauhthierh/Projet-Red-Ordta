package library

import "testing"

func TestOrCombat3D(t *testing.T) {
	for _, genre := range []TypeMonstre{TypeCorbeau, TypeSanglier, TypeLoup, TypeTroll} {
		p := joueurCombatTest()
		p.Attaque = 1000
		// Une bourse n'occupe pas de case : récompense même avec le sac plein.
		p.CapaciteInventaire = p.TotalInventaire()
		avant := p.Argent
		c, err := NouveauCombat3D(&p, ModeDuel, genre)
		if err != nil {
			t.Fatal(err)
		}
		forcerJoueurCommence(c)
		// L'or est tiré au hasard : on vérifie qu'il reste dans la fourchette du monstre.
		orMin := c.Ennemis[0].Monstre.OrMin
		orMax := c.Ennemis[0].Monstre.OrMax
		r := c.Attaquer(0)
		if !r.Victoire || r.OrGagne < orMin || r.OrGagne > orMax {
			t.Fatalf("récompense hors fourchette pour %s (%d à %d) : %+v", genre, orMin, orMax, r)
		}
		if c.OrTotal != r.OrGagne || p.Argent != avant+r.OrGagne {
			t.Fatalf("or mal ajouté pour %s : %+v", genre, r)
		}
		gain := r.OrGagne
		c.Attaquer(0)
		c.Quitter3D()
		c.Quitter3D()
		if p.Argent != avant+gain {
			t.Fatal("or perdu ou attribué deux fois")
		}
	}
}

func TestPasOrApresAbandon(t *testing.T) {
	p := joueurCombatTest()
	avant := p.Argent
	c, _ := NouveauCombat3D(&p, ModeDuel, TypeLoup)
	c.Quitter3D()
	if p.Argent != avant || c.OrTotal != 0 {
		t.Fatal("l'abandon donne de l'or")
	}
}
