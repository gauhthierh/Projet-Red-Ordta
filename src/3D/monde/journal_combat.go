package monde

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"ordta/library"
	"strings"
)

// Le journal conserve le combat entier ; les boutons parcourent les événements.
type JournalCombat struct {
	Panneau *gui.Panel
	texte   *gui.Label
	entrees []string
	index   int
	phase   string
}

func NouveauJournalCombat(scene *core.Node) *JournalCombat {
	j := &JournalCombat{Panneau: gui.NewPanel(480, 310)}
	j.Panneau.SetPosition(30, 290)
	j.Panneau.SetColor4(&math32.Color4{R: .03, G: .03, B: .05, A: .9})
	j.texte = gui.NewLabel("")
	j.texte.SetPosition(15, 15)
	j.Panneau.Add(j.texte)
	for index, nom := range []string{"PRÉCÉDENT", "SUIVANT"} {
		pas := index*2 - 1
		b := gui.NewButton(nom)
		b.SetPosition(15+float32(index)*225, 265)
		b.SetSize(210, 30)
		b.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
			if j.index+pas >= 0 && j.index+pas < len(j.entrees) {
				j.index += pas
				j.afficher()
			}
		})
		j.Panneau.Add(b)
	}
	j.Panneau.SetVisible(false)
	scene.Add(j.Panneau)
	return j
}

func (j *JournalCombat) afficher() {
	if len(j.entrees) == 0 {
		j.texte.SetText("")
		return
	}
	j.texte.SetText(fmt.Sprintf("JOURNAL — %d / %d\n\n%s", j.index+1, len(j.entrees), j.entrees[j.index]))
}

func (j *JournalCombat) Ajouter(message string) {
	// Une page courte par événement : aucune récompense ne disparaît du journal.
	var lignes []string
	for _, ligne := range strings.Split(message, "\n") {
		courante := ""
		for _, mot := range strings.Fields(ligne) {
			if len([]rune(courante+mot)) > 53 {
				lignes = append(lignes, courante)
				courante = ""
			}
			courante += mot + " "
		}
		lignes = append(lignes, courante)
	}
	for len(lignes) > 0 {
		n := min(10, len(lignes))
		j.entrees = append(j.entrees, strings.Join(lignes[:n], "\n"))
		lignes = lignes[n:]
	}
	j.index = len(j.entrees) - 1
	j.afficher()
}

func (j *JournalCombat) Reinitialiser() { j.entrees = nil; j.phase = ""; j.index = 0; j.afficher() }

func (j *JournalCombat) Etat(c *library.CombatArene) {
	cle := fmt.Sprintf("Vague %d — tour %d — %s", c.NumeroVague, c.Tour, c.Phase)
	if cle != j.phase {
		j.phase = cle
		j.Ajouter(cle)
	}
}

func texteResultatCombat(action string, r library.ResultatAction) string {
	t := action + "\n" + r.Message
	if r.Degats > 0 {
		t += fmt.Sprintf("\nDégâts : %d | %s : %d → %d PV", r.Degats, r.Cible, r.PVAvant, r.PVApres)
	}
	if r.Soin > 0 {
		t += fmt.Sprintf("\nSoin : +%d PV (%d → %d)", r.Soin, r.PVAvant, r.PVApres)
	}
	if r.ManaAvant != r.ManaApres {
		t += fmt.Sprintf("\nMana : %+d (%d → %d)", r.ManaApres-r.ManaAvant, r.ManaAvant, r.ManaApres)
	}
	if r.CibleVaincue {
		t += "\n" + r.Cible + " vaincu."
	}
	if r.ExperienceGagnee != 0 || r.OrGagne != 0 {
		t += fmt.Sprintf("\nRécompenses : +%d XP, +%d or", r.ExperienceGagnee, r.OrGagne)
	}
	if r.NiveauxGagnes > 0 {
		t += fmt.Sprintf("\nNiveaux gagnés : %d", r.NiveauxGagnes)
	}
	return t
}
