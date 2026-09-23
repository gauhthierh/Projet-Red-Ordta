package monde

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"ordta/library"
)

type ChoixCombat struct {
	Panneau                            *gui.Panel
	Message                            *gui.Label
	Entrainement, Arene, Duel, Quitter *gui.Button
	Creatures                          map[library.TypeMonstre]*gui.Button
}

const aideChoixCombat = "CHOISIR UN COMBAT\nEntraînement : aucun gain, PV / mana / objets restaurés à la sortie.\nArène : vagues actuelles. Duel : un seul monstre pour les matériaux.\nVous pouvez quitter à tout moment. Prévoyez des places pour le butin."

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
		b := bouton(fmt.Sprintf("%s — %d PV", ennemi.Monstre.Nom, ennemi.Monstre.PVMax), 25+float32(index%2)*410, 220+float32(index/2)*100)
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

func (c *ChoixCombat) Ouvrir() {
	c.Message.SetText(aideChoixCombat)
	for _, b := range c.Creatures {
		b.SetVisible(false)
	}
	c.Panneau.SetVisible(true)
}
