package monde

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"ordta/library"
)

// Ce panneau ne s'ouvre que lorsqu'aucune sauvegarde n'existe.
type CreationPersonnage struct {
	Panneau         *gui.Panel
	Nom             *gui.Edit
	Message         *gui.Label
	Valider, Retour *gui.Button
	Classe          string
}

// NouvelleCreationPersonnage : Prépare la saisie du nom et le choix de classe ; la validation est branchée dans monde.go.
func NouvelleCreationPersonnage(scene *core.Node, largeur, hauteur float32) *CreationPersonnage {
	c := &CreationPersonnage{Panneau: gui.NewPanel(760, 520), Classe: "Humain"}
	c.Panneau.SetPosition((largeur-760)/2, (hauteur-520)/2)
	c.Panneau.SetColor4(&math32.Color4{R: .03, G: .04, B: .07, A: 1})
	titre := gui.NewLabel("CRÉER VOTRE PERSONNAGE")
	titre.SetPosition(30, 25)
	titre.SetFontSize(24)
	c.Panneau.Add(titre)
	c.Nom = gui.NewEdit(680, "Nom : lettres uniquement")
	c.Nom.SetPosition(30, 85)
	c.Panneau.Add(c.Nom)
	statistiques := gui.NewLabel("")
	statistiques.SetPosition(30, 230)
	c.Panneau.Add(statistiques)
	choisir := func(nom string) {
		c.Classe = nom
		classe, _ := library.ClasseDepuisNom3D(nom)
		statistiques.SetText(fmt.Sprintf("%s\nPV de départ : %d / %d — Mana : %d — Attaque : %d\nNiveau 1 — 100 or — 3 potions — 10 places", nom, classe.PvMax/2, classe.PvMax, classe.ManaMax, classe.Attaque))
	}
	for index, nom := range []string{"Humain", "Elfe", "Nain"} {
		classe := nom
		bouton := gui.NewButton(nom)
		bouton.SetSize(210, 50)
		bouton.SetPosition(30+float32(index)*235, 150)
		bouton.Subscribe(gui.OnClick, func(_ string, _ interface{}) { choisir(classe) })
		c.Panneau.Add(bouton)
	}
	choisir("Humain")
	c.Message = gui.NewLabel("")
	c.Message.SetPosition(30, 350)
	c.Panneau.Add(c.Message)
	c.Valider = gui.NewButton("COMMENCER")
	c.Valider.SetSize(260, 55)
	c.Valider.SetPosition(30, 420)
	c.Retour = gui.NewButton("RETOUR")
	c.Retour.SetSize(260, 55)
	c.Retour.SetPosition(450, 420)
	c.Panneau.Add(c.Valider)
	c.Panneau.Add(c.Retour)
	c.Panneau.SetVisible(false)
	scene.Add(c.Panneau)
	return c
}
