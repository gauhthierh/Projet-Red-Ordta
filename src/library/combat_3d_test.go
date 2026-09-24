package library

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func joueurCombatTest() Character {
	p := NouveauPersonnage3D("Test", "Humain")
	p.PvActuel = p.PvMaxTotal
	return p
}

// forcerJoueurCommence : L'initiative est tirée au hasard ; les tests qui supposent
// que le joueur joue en premier le fixent eux-mêmes pour ne pas dépendre du tirage.
func forcerJoueurCommence(c *CombatArene) {
	c.joueurCommence = true
	c.Phase = PhaseTourJoueur
}

func TestModesCombat3D(t *testing.T) {
	for _, mode := range []ModeCombat{ModeEntrainement, ModeArene, ModeDuel} {
		p := joueurCombatTest()
		c, err := NouveauCombat3D(&p, mode, TypeLoup)
		if err != nil {
			t.Fatal(err)
		}
		// Selon l'initiative, l'un ou l'autre camp commence, toujours au tour 1.
		if (c.Phase != PhaseTourJoueur && c.Phase != PhaseTourMonstres) || c.Tour != 1 {
			t.Fatal("le combat doit commencer au tour 1 par l'un des deux camps")
		}
		if mode == ModeArene && len(c.Vagues) != 4 {
			t.Fatal("vagues de l'arène modifiées")
		}
		if mode != ModeArene && (len(c.Vagues) != 1 || len(c.Ennemis) != 1) {
			t.Fatal("un seul adversaire attendu")
		}
	}
	if _, err := NouveauCombat3D(nil, ModeDuel, TypeLoup); err == nil {
		t.Fatal("joueur absent accepté")
	}
	p := joueurCombatTest()
	if _, err := NouveauCombat3D(&p, ModeDuel, TypeMonstre("inconnu")); err == nil {
		t.Fatal("monstre inconnu accepté")
	}
}

func TestDuelToursVictoireEtButinUnique(t *testing.T) {
	p := joueurCombatTest()
	c, _ := NouveauCombat3D(&p, ModeDuel, TypeLoup)
	forcerJoueurCommence(c)
	c.tirageButin = func(int) int { return 0 }
	r := c.Attaquer(0)
	if !r.Reussite || c.Phase != PhaseTourMonstres {
		t.Fatal("l'attaque doit céder le tour")
	}
	if c.Attaquer(0).Reussite {
		t.Fatal("double action du joueur acceptée")
	}
	c.ProchaineActionMonstre()
	if c.Phase != PhaseTourJoueur || c.Tour != 2 {
		t.Fatal("tour suivant incorrect")
	}
	p.Attaque = 1000
	r = c.Attaquer(0)
	if !r.Victoire || c.Phase != PhaseVictoire || p.Inventaire[ItemFourrureDeLoup] != 1 {
		t.Fatalf("victoire ou butin incorrect : %+v", r)
	}
	if c.Attaquer(0).Reussite || p.Inventaire[ItemFourrureDeLoup] != 1 {
		t.Fatal("butin obtenu deux fois")
	}
	if c.CommencerVagueSuivante() == nil {
		t.Fatal("le duel doit s'arrêter après un adversaire")
	}
	c.Quitter3D()
	if p.Inventaire[ItemFourrureDeLoup] != 1 {
		t.Fatal("butin perdu en quittant")
	}
}

func TestButinsPourcentagesEtRecettes(t *testing.T) {
	materiaux := map[string]bool{}
	for _, genre := range []TypeMonstre{TypeLoup, TypeTroll, TypeSanglier, TypeCorbeau} {
		butin := ButinPossible3D(genre)
		materiaux[butin.Objet] = true
		p := joueurCombatTest()
		p.CapaciteInventaire = 200
		c, _ := NouveauCombat3D(&p, ModeDuel, genre)
		for tirage := 0; tirage < 100; tirage++ {
			c.tirageButin = func(int) int { return tirage }
			c.donnerButin(genre)
		}
		if p.Inventaire[butin.Objet] != butin.Pourcentage {
			t.Fatalf("probabilité incorrecte pour %s", genre)
		}
		p.CapaciteInventaire = p.TotalInventaire()
		c.tirageButin = func(int) int { return 0 }
		c.donnerButin(genre)
		if p.TotalInventaire() != p.CapaciteInventaire {
			t.Fatal("inventaire dépassé")
		}
	}
	for _, armure := range Armurerie {
		for materiau := range armure.Materiaux {
			if !materiaux[materiau] {
				t.Fatalf("matériau impossible à récolter : %s", materiau)
			}
		}
	}
}

func TestAbandonStoppeActionsEtPoison(t *testing.T) {
	for _, phase := range []PhaseCombat{PhaseTourJoueur, PhaseTourMonstres, PhaseEntreVagues} {
		p := joueurCombatTest()
		c, _ := NouveauCombat3D(&p, ModeDuel, TypeSanglier)
		c.Phase = phase
		c.poisonSecondes = 3
		c.Quitter3D()
		pv := p.PvActuel
		if c.Attaquer(0).Reussite || c.ProchaineActionMonstre().Reussite || len(c.MettreAJourEffets(time.Minute)) != 0 {
			t.Fatal("action après abandon")
		}
		if p.PvActuel != pv || len(c.Butins) > 0 || c.ExperienceTotale != 0 {
			t.Fatal("abandon modifie PV ou récompenses")
		}
	}
}

func TestModelesDuelDisponibles(t *testing.T) {
	for _, genre := range []TypeMonstre{TypeLoup, TypeTroll, TypeSanglier, TypeCorbeau} {
		ennemi, _ := EnnemiDuel3D(genre)
		data, err := os.ReadFile(filepath.Join("../../assets/models/monstres", ennemi.Modele3D))
		if err != nil {
			t.Fatal(err)
		}
		var modele struct {
			Parties []struct {
				Pieces []struct {
					Forme   string
					Taille  [3]float64
					Texture string
				}
			}
		}
		if err := json.Unmarshal(data, &modele); err != nil {
			t.Fatal(err)
		}
		if len(modele.Parties) == 0 {
			t.Fatal("modèle vide")
		}
		for _, partie := range modele.Parties {
			for _, p := range partie.Pieces {
				if p.Forme != "boite" && p.Forme != "ellipsoide" {
					t.Fatal("forme inconnue")
				}
				for _, s := range p.Taille {
					if s <= 0 {
						t.Fatal("dimension incorrecte")
					}
				}
				if p.Texture != "" {
					if _, err := os.Stat(filepath.Join("../../assets/models/monstres/textures", p.Texture+".png")); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}
