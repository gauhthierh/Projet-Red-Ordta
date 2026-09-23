package monde

import (
	"github.com/g3n/engine/core"
	"ordta/library"
	"strings"
	"testing"
)

func TestCombatAnimationRestaureLesPoses(t *testing.T) {
	for _, action := range []string{"attaque", "gobelin", "gobelin_cuirasse", "loup", "sanglier", "corbeau", "troll", "chaman", "soin_monstre", "defense", "objet",
		library.SortCoupDePoing, library.SortLameDuDestin, library.SortBouleDeFeu, library.SortEclateDuGardien, library.SortFlecheDeLumiere, library.SortFoudreCeleste, library.SortSoinDuCoeur, library.SortBouclier, library.SortDevotion, library.SortDernierEspoir, library.SortJugementDesGeants} {
		t.Run(action, func(t *testing.T) {
			scene, acteur, cible, bras := core.NewNode(), core.NewNode(), core.NewNode(), core.NewNode()
			acteur.SetPosition(1, 2, 3)
			cible.SetPosition(8, 2, 3)
			bras.SetRotation(.2, .3, .4)
			position, rotation := acteur.Position(), bras.Rotation()
			a := NouvelleAnimationCombat(scene)
			termine := 0
			a.Demarrer(action, acteur, cible, map[string]*core.Node{"bras_droit": bras}, func() { termine++ })
			a.MettreAJour(.6)
			if !a.Active {
				t.Fatal("animation interrompue trop tôt")
			}
			a.MettreAJour(1)
			a.MettreAJour(1)
			if a.Active || termine != 1 || acteur.Position() != position || bras.Rotation() != rotation {
				t.Fatal("pose ou fin incorrecte")
			}
			a.Demarrer(action, acteur, cible, nil, func() { t.Error("callback appelé après annulation") })
			a.MettreAJour(.5)
			a.Arreter()
			if acteur.Position() != position {
				t.Fatal("la sortie doit restaurer la position")
			}
		})
	}
}

func TestTexteCombatContientLesChiffres(t *testing.T) {
	r := library.ResultatAction{Message: "Impact", Cible: "Gobelin", Degats: 12, PVAvant: 40, PVApres: 28, ManaAvant: 50, ManaApres: 30, ExperienceGagnee: 10, OrGagne: 4}
	texte := texteResultatCombat("Sort", r)
	for _, attendu := range []string{"12", "40 → 28", "-20 (50 → 30)", "+10 XP, +4 or"} {
		if !strings.Contains(texte, attendu) {
			t.Errorf("détail absent : %s", attendu)
		}
	}
}

func TestAttaqueBasiqueBrasVersAvant(t *testing.T) {
	for _, action := range []string{"Attaque Basique", "attaque"} {
		bras := core.NewNode()
		a := NouvelleAnimationCombat(core.NewNode())
		a.Demarrer(action, core.NewNode(), core.NewNode(), map[string]*core.Node{"bras_droit": bras}, nil)
		a.MettreAJour(.625)
		if bras.Rotation().X >= 0 {
			t.Errorf("%s : le bras doit tourner vers l'avant (-X)", action)
		}
		a.MettreAJour(.625)
		if bras.Rotation().X != 0 {
			t.Error("le bras doit retrouver sa pose initiale")
		}
	}
}
