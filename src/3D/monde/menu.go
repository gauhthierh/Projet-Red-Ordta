package monde

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

// Le même panneau sert à l'accueil puis à la pause.
type MenuJeu struct {
	Panneau         *gui.Panel
	Titre           *gui.Label
	BoutonReprendre *gui.Button
	BoutonDemarrer  *gui.Button
	BoutonQuitter   *gui.Button
	Ouvert          bool
	Demarre         bool
}

// NouveauMenuJeu : Prépare le menu commun à l'accueil et à la pause.
func NouveauMenuJeu(scene *core.Node, largeur, hauteur float32) *MenuJeu {
	fond := gui.NewPanel(largeur, hauteur)
	fond.SetColor4(&math32.Color4{R: .025, G: .035, B: .055, A: .94})
	titre := gui.NewLabel("Ordta")
	titre.SetFontSize(32)
	titre.SetPosition(largeur/2-40, hauteur/2-145)
	fond.Add(titre)
	demarrer := gui.NewButton("DÉMARRER")
	demarrer.SetSize(300, 60)
	demarrer.SetPosition(largeur/2-150, hauteur/2-40)
	fond.Add(demarrer)
	reprendre := gui.NewButton("REPRENDRE [ÉCHAP]")
	reprendre.SetSize(300, 60)
	reprendre.SetPosition(largeur/2-150, hauteur/2-40)
	fond.Add(reprendre)
	reprendre.SetVisible(false)
	quitter := gui.NewButton("QUITTER LE JEU")
	quitter.SetSize(300, 60)
	quitter.SetPosition(largeur/2-150, hauteur/2+45)
	fond.Add(quitter)
	artistes := gui.NewButton("QUI SONT-ILS ?")
	artistes.SetSize(300, 50)
	artistes.SetPosition(largeur/2-150, hauteur/2+125)
	fond.Add(artistes)
	reponse := gui.NewLabel("ABBA et Steven Spielberg")
	reponse.SetPosition(largeur/2-80, hauteur/2+195)
	reponse.SetVisible(false)
	fond.Add(reponse)
	artistes.Subscribe(gui.OnClick, func(_ string, _ interface{}) { reponse.SetVisible(true) })
	scene.Add(fond)
	return &MenuJeu{Panneau: fond, Titre: titre, BoutonDemarrer: demarrer, BoutonReprendre: reprendre, BoutonQuitter: quitter, Ouvert: true}
}

// Reprendre : Masque le menu et marque la partie comme démarrée.
func (m *MenuJeu) Reprendre() {
	m.Demarre = true
	m.Ouvert = false
	m.Panneau.SetVisible(false)
}

// Pause : Affiche le menu de reprise ; la boucle principale suspend alors la simulation.
func (m *MenuJeu) Pause() {
	m.Ouvert = true
	m.Titre.SetText("PAUSE")
	m.BoutonDemarrer.SetVisible(false)
	m.BoutonReprendre.SetVisible(true)
	m.Panneau.SetVisible(true)
}
