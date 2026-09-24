package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ordta/library"
)

// Crée un joueur de test en pleine santé //
func joueurPourCombat() library.Character {
	p := library.NouveauPersonnage3D("Test", "Humain")
	p.PvActuel = p.PvMaxTotal
	return p
}

// L'initiative est tirée au hasard et ne peut pas être fixée depuis ce package //
// On recrée donc le combat jusqu'à ce que le joueur commence //
func combatOuJoueurCommence(t *testing.T, p *library.Character, mode library.ModeCombat, genre library.TypeMonstre) *library.CombatArene {
	t.Helper()
	for essai := 0; essai < 100; essai++ {
		c, err := library.NouveauCombat3D(p, mode, genre)
		if err != nil {
			t.Fatal(err)
		}
		if c.Phase == library.PhaseTourJoueur {
			return c
		}
	}
	t.Fatal("le joueur n'a jamais commencé en 100 essais")
	return nil
}

func TestModesCombat3D(t *testing.T) {
	for _, mode := range []library.ModeCombat{library.ModeEntrainement, library.ModeArene, library.ModeDuel} {
		p := joueurPourCombat()
		c, err := library.NouveauCombat3D(&p, mode, library.TypeLoup)
		if err != nil {
			t.Fatal(err)
		}
		// Selon l'initiative, l'un ou l'autre camp commence, toujours au tour 1 //
		if (c.Phase != library.PhaseTourJoueur && c.Phase != library.PhaseTourMonstres) || c.Tour != 1 {
			t.Fatal("le combat doit commencer au tour 1 par l'un des deux camps")
		}
		if mode == library.ModeArene && len(c.Vagues) != 4 {
			t.Fatal("vagues de l'arène modifiées")
		}
		if mode != library.ModeArene && (len(c.Vagues) != 1 || len(c.Ennemis) != 1) {
			t.Fatal("un seul adversaire attendu")
		}
	}
	if _, err := library.NouveauCombat3D(nil, library.ModeDuel, library.TypeLoup); err == nil {
		t.Fatal("joueur absent accepté")
	}
	p := joueurPourCombat()
	if _, err := library.NouveauCombat3D(&p, library.ModeDuel, library.TypeMonstre("inconnu")); err == nil {
		t.Fatal("monstre inconnu accepté")
	}
}

func TestDuelToursVictoireEtButinUnique(t *testing.T) {
	p := joueurPourCombat()
	c := combatOuJoueurCommence(t, &p, library.ModeDuel, library.TypeLoup)
	r := c.Attaquer(0)
	if !r.Reussite || c.Phase != library.PhaseTourMonstres {
		t.Fatal("l'attaque doit céder le tour")
	}
	if c.Attaquer(0).Reussite {
		t.Fatal("double action du joueur acceptée")
	}
	c.ProchaineActionMonstre()
	if c.Phase != library.PhaseTourJoueur || c.Tour != 2 {
		t.Fatal("tour suivant incorrect")
	}
	p.Attaque = 1000
	r = c.Attaquer(0)
	if !r.Victoire || c.Phase != library.PhaseVictoire {
		t.Fatalf("victoire incorrecte : %+v", r)
	}
	// Le butin est tiré au hasard : on vérifie seulement qu'il n'est jamais donné deux fois //
	fourrures := p.Inventaire[library.ItemFourrureDeLoup]
	if fourrures > 1 {
		t.Fatal("plus d'une fourrure pour un seul loup")
	}
	if c.Attaquer(0).Reussite || p.Inventaire[library.ItemFourrureDeLoup] != fourrures {
		t.Fatal("butin obtenu deux fois")
	}
	if c.CommencerVagueSuivante() == nil {
		t.Fatal("le duel doit s'arrêter après un adversaire")
	}
	c.Quitter3D()
	if p.Inventaire[library.ItemFourrureDeLoup] != fourrures {
		t.Fatal("butin perdu en quittant")
	}
}

func TestButinsRecettes(t *testing.T) {
	materiaux := map[string]bool{}
	for _, genre := range []library.TypeMonstre{library.TypeLoup, library.TypeTroll, library.TypeSanglier, library.TypeCorbeau} {
		butin := library.ButinPossible3D(genre)
		if butin.Pourcentage <= 0 || butin.Pourcentage > 100 {
			t.Fatalf("pourcentage de butin incorrect pour %s", genre)
		}
		materiaux[butin.Objet] = true
	}
	for _, armure := range library.Armurerie {
		for materiau := range armure.Materiaux {
			if !materiaux[materiau] {
				t.Fatalf("matériau impossible à récolter : %s", materiau)
			}
		}
	}
}

func TestAbandonStoppeActionsEtPoison(t *testing.T) {
	for _, phase := range []library.PhaseCombat{library.PhaseTourJoueur, library.PhaseTourMonstres, library.PhaseEntreVagues} {
		p := joueurPourCombat()
		p.Inventaire[library.ItemPotionDePoison] = 1
		c := combatOuJoueurCommence(t, &p, library.ModeDuel, library.TypeSanglier)
		// Le poison est lancé normalement, puisque son compteur est privé //
		if !c.UtiliserObjet(library.ItemPotionDePoison).Reussite {
			t.Fatal("la potion de poison doit pouvoir être bue")
		}
		c.Phase = phase
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
	for _, genre := range []library.TypeMonstre{library.TypeLoup, library.TypeTroll, library.TypeSanglier, library.TypeCorbeau} {
		ennemi, _ := library.EnnemiDuel3D(genre)
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
			for _, piece := range partie.Pieces {
				if piece.Forme != "boite" && piece.Forme != "ellipsoide" {
					t.Fatal("forme inconnue")
				}
				for _, s := range piece.Taille {
					if s <= 0 {
						t.Fatal("dimension incorrecte")
					}
				}
				if piece.Texture != "" {
					if _, err := os.Stat(filepath.Join("../../assets/models/monstres/textures", piece.Texture+".png")); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}
