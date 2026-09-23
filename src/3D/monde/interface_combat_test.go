package monde

import (
	"github.com/g3n/engine/core"
	"ordta/library"
	"testing"
)

// Les options de sort ne doivent pas survivre au passage au tour adverse.
func TestCommandesSuiventPhaseCombat(t *testing.T) {
	i := NouvelleInterfaceCombat(core.NewNode())
	c := &library.CombatArene{Phase: library.PhaseTourJoueur}
	i.ActualiserVisibilite(c, false, false)
	if !i.BoutonAttaquer.Visible() {
		t.Fatal("commandes absentes au tour joueur")
	}
	i.AfficherChoix("sort", []string{library.SortCoupDePoing})
	c.Phase = library.PhaseTourMonstres
	i.ActualiserVisibilite(c, false, false)
	if i.BoutonAttaquer.Visible() || len(i.OptionsChoix) != 0 {
		t.Fatal("anciens choix encore actifs")
	}
	c.Phase = library.PhaseVictoire
	i.ActualiserVisibilite(c, false, false)
	if !i.BoutonRejouer.Visible() {
		t.Fatal("rejouer absent après victoire")
	}
	i.ActualiserVisibilite(nil, false, false)
	if i.Panneau.Visible() {
		t.Fatal("commandes visibles hors combat")
	}
}
