package monde

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"ordta/library"
)

// ChoixCombat : Rassemble les boutons qui précèdent le démarrage d'un combat.
type ChoixCombat struct {
	Panneau                            *gui.Panel
	Message                            *gui.Label
	Entrainement, Arene, Duel, Quitter *gui.Button
	Creatures                          map[library.TypeMonstre]*gui.Button
}

const aideChoixCombat = "CHOISIR UN COMBAT\nEntraînement : gobelin, expérience et or gagnés.\nArène : vagues. Duel : un monstre pour les matériaux.\nEnnemis adaptés à votre niveau. Vous pouvez quitter à tout moment."

// NouveauChoixCombat : Construit le choix des modes et des créatures ; le démarrage est raccordé dans monde.go.
func NouveauChoixCombat(scene *core.Node) *ChoixCombat {
	c := &ChoixCombat{Panneau: gui.NewPanel(860, 550), Creatures: map[library.TypeMonstre]*gui.Button{}}
	c.Panneau.SetPosition(530, 240)
	c.Panneau.SetColor4(&math32.Color4{R: .04, G: .04, B: .07, A: .97})
	c.Message = gui.NewLabel(aideChoixCombat)
	c.Message.SetPosition(25, 25)
	c.Panneau.Add(c.Message)
	bouton := func(nom string, x, y float32) *gui.Button {
		b := gui.NewButton(nom)
		b.SetPosition(x, y)
		b.SetSize(250, 50)
		c.Panneau.Add(b)
		return b
	}
	c.Entrainement = bouton("ENTRAÎNEMENT", 25, 140)
	c.Arene = bouton("ARÈNE À VAGUES", 305, 140)
	c.Duel = bouton("DUEL : CHOISIR UN MONSTRE", 585, 140)
	for index, genre := range []library.TypeMonstre{library.TypeCorbeau, library.TypeSanglier, library.TypeLoup, library.TypeTroll} {
		ennemi, _ := library.EnnemiDuel3D(genre)
		butin := library.ButinPossible3D(genre)
		b := bouton(fmt.Sprintf("%s — %d PV", ennemi.Monstre.Nom, ennemi.Monstre.PvMax), 25+float32(index%2)*410, 220+float32(index/2)*100)
		b.SetSize(385, 50)
		b.SetVisible(false)
		c.Creatures[genre] = b
		texte := gui.NewLabel(fmt.Sprintf("%s : %d %%", butin.Objet, butin.Pourcentage))
		texte.SetPosition(25+float32(index%2)*410, 275+float32(index/2)*100)
		c.Panneau.Add(texte)
	}
	c.Duel.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		for _, b := range c.Creatures {
			b.SetVisible(true)
		}
	})
	c.Quitter = bouton("QUITTER L'ARÈNE", 305, 460)
	c.Panneau.SetVisible(false)
	scene.Add(c.Panneau)
	return c
}

// Ouvrir : Réinitialise le panneau de choix avant une nouvelle session de combat.
func (c *ChoixCombat) Ouvrir(niveau int) {
	c.Message.SetText(aideChoixCombat)
	for genre, b := range c.Creatures {
		ennemi, _ := library.EnnemiDuel3D(genre)
		m := library.AdapterEnnemiNiveau3D(ennemi, niveau).Monstre
		b.Label.SetText(fmt.Sprintf("%s niv. %d — %d PV — %d–%d or", m.Nom, m.Niveau, m.PvMax, m.OrMin, m.OrMax))
		b.SetVisible(false)
	}
	c.Panneau.SetVisible(true)
}
