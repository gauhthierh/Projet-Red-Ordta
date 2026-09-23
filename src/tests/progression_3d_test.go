package tests

import (
	"ordta/library"
	"testing"
)

// Les sept créatures gardent leurs bases, mais suivent le niveau du joueur.
func TestProgressionToutesCreatures3D(t *testing.T) {
	ennemis := []library.EnnemiCombat{
		library.NouveauGobelinCombat(), library.NouveauGobelinCuirasseCombat(),
		library.NouveauChamanCombat(), library.NouveauLoupCombat(), library.NouveauTrollCombat(),
	}
	for _, genre := range []library.TypeMonstre{library.TypeSanglier, library.TypeCorbeau} {
		e, err := library.EnnemiDuel3D(genre)
		if err != nil {
			t.Fatal(err)
		}
		ennemis = append(ennemis, e)
	}
	for _, base := range ennemis {
		for _, niveau := range []int{1, 3, 10} {
			e := library.AdapterEnnemiNiveau3D(base, niveau)
			m := e.Monstre
			bonus := niveau - 1
			if m.Niveau != niveau || m.PvMax != base.Monstre.PvMax+library.GainPvGobelin*bonus ||
				m.PvActuel != m.PvMax || m.Attaque != base.Monstre.Attaque+library.GainAttaqueGobelin*bonus ||
				m.ExperienceDonnee != base.Monstre.ExperienceDonnee+library.GainExperienceGobelin*bonus ||
				m.OrMin != base.Monstre.OrMin+library.GainOrMinGobelin*bonus ||
				m.OrMax != base.Monstre.OrMax+library.GainOrMaxGobelin*bonus || m.OrMin <= 0 || m.OrMax < m.OrMin {
				t.Fatalf("%s niveau %d : %+v", e.Type, niveau, m)
			}
			if library.AdapterEnnemiNiveau3D(e, niveau).Monstre != m {
				t.Fatal("les bonus sont appliqués deux fois")
			}
			if base.Type == library.TypeGobelin && m != library.InitGoblin(niveau) {
				t.Fatal("gobelin différent du backend")
			}
		}
	}
}

// Une élimination rapporte une seule fois, même en entraînement.
func TestOrTousModesEtSortie3D(t *testing.T) {
	for _, mode := range []library.ModeCombat{library.ModeEntrainement, library.ModeArene, library.ModeDuel} {
		p := library.NouveauPersonnage3D("Test", "Humain")
		p.Niveau, p.Attaque, p.PvActuel = 3, 10000, 100
		c, err := library.NouveauCombat3D(&p, mode, library.TypeTroll)
		if err != nil {
			t.Fatal(err)
		}
		m := c.Ennemis[0].Monstre
		if m.Niveau != p.Niveau {
			t.Fatal("niveau non appliqué à la session")
		}
		c.Phase = library.PhaseTourJoueur
		avant := p.Argent
		r := c.Attaquer(0)
		if !r.CibleVaincue || r.OrGagne < m.OrMin || r.OrGagne > m.OrMax ||
			p.Argent != avant+r.OrGagne || c.OrTotal != r.OrGagne || r.ExperienceGagnee != m.ExperienceDonnee {
			t.Fatalf("%s : %+v", mode, r)
		}
		apres := p.Argent
		if c.Attaquer(0).Reussite {
			t.Fatal("attaque sur une cible morte")
		}
		c.Quitter3D()
		if p.Argent != apres {
			t.Fatal("gain perdu ou doublé à la sortie")
		}
	}
	p := library.NouveauPersonnage3D("Test", "Humain")
	c, _ := library.NouveauCombatArene(&p)
	avant := p.Argent
	c.Quitter3D()
	if p.Argent != avant {
		t.Fatal("abandon récompensé sans élimination")
	}
}

func TestNiveauFigePuisActualiseVagueSuivante(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	c, _ := library.NouveauCombatArene(&p)
	p.Niveau = 4
	if c.Ennemis[0].Monstre.Niveau != 1 {
		t.Fatal("ennemi modifié en cours de vague")
	}
	c.Phase = library.PhaseEntreVagues
	if err := c.CommencerVagueSuivante(); err != nil {
		t.Fatal(err)
	}
	for _, e := range c.Ennemis {
		if e.Monstre.Niveau != 4 {
			t.Fatal("nouvelle vague non adaptée")
		}
	}
	if c.Vagues[1][0].Monstre.Niveau != 1 {
		t.Fatal("définition de vague modifiée")
	}
}

func TestAchatsInutilesEtArmureEnCombat3D(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	p.Niveau, p.Argent = 10, 1000
	for index, article := range library.BoutiqueMarchand3D() {
		if article.Nom != library.ItemLivreLameDuDestin {
			continue
		}
		if !p.AcheterMarchand3D(index).Reussite {
			t.Fatal("premier livre refusé")
		}
		or := p.Argent
		if p.AcheterMarchand3D(index).Reussite || p.Argent != or {
			t.Fatal("livre possédé racheté")
		}
		p.UtiliserObjet3D(article.Nom)
		if p.AcheterMarchand3D(index).Reussite {
			t.Fatal("sort appris racheté")
		}
	}
	p.Inventaire[library.ItemChapeauAventurier] = 1
	c, _ := library.NouveauCombatArene(&p)
	c.Phase = library.PhaseTourJoueur
	if c.UtiliserObjet(library.ItemChapeauAventurier).Reussite ||
		p.Equipement.Tete.Nom != "" || p.Inventaire[library.ItemChapeauAventurier] != 1 ||
		c.Phase != library.PhaseTourJoueur {
		t.Fatal("équipement modifié ou tour consommé")
	}
}
