package monde

import (
	"fmt"
	"ordta/3D/personnages"
	"ordta/library"
	"sort"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type TypePNJ string

const (
	PNJMarchand   TypePNJ = "marchand"
	PNJForgeron   TypePNJ = "forgeron"
	rayonDialogue         = float32(4.5)
)

type PNJ struct {
	Type     TypePNJ
	Nom      string
	Position math32.Vector3
	Modele   *personnages.Personnage
}

type InterfaceInteraction struct {
	Panneau *gui.Panel
	Texte   *gui.Label
}

type InterfaceCommerce struct {
	Panneau *gui.Panel
	Titre   *gui.Label
	Infos   *gui.Label
	Message *gui.Label
	Boutons []*gui.Button
	Fermer  *gui.Button

	Ouvert bool
	Mode   TypePNJ
	Joueur *library.Character
}

func AjouterPNJMarche(scene *core.Node) []*PNJ {
	definitions := []struct {
		typePNJ  TypePNJ
		nom      string
		position math32.Vector3
	}{
		{PNJMarchand, "Marchand", math32.Vector3{X: 35, Y: 12.5, Z: .265}},
		{PNJForgeron, "Forgeron", math32.Vector3{X: 48, Y: 8.5, Z: .265}},
	}

	pnjs := make([]*PNJ, 0, len(definitions))
	for _, definition := range definitions {
		modele := personnages.NouveauAvecTenue(string(definition.typePNJ))
		noeud := modele.Noeud()
		noeud.SetPosition(definition.position.X, definition.position.Y, definition.position.Z)
		noeud.SetRotationZ(0)
		scene.Add(noeud)
		pnjs = append(pnjs, &PNJ{
			Type: definition.typePNJ, Nom: definition.nom,
			Position: definition.position, Modele: modele,
		})
	}
	return pnjs
}

func PNJLePlusProche(position math32.Vector3, pnjs []*PNJ) *PNJ {
	var proche *PNJ
	meilleureDistance := rayonDialogue * rayonDialogue
	for _, pnj := range pnjs {
		differenceX := position.X - pnj.Position.X
		differenceY := position.Y - pnj.Position.Y
		distanceCarree := differenceX*differenceX + differenceY*differenceY
		if distanceCarree <= meilleureDistance {
			proche = pnj
			meilleureDistance = distanceCarree
		}
	}
	return proche
}

func NouvelleInterfaceInteraction(scene *core.Node) *InterfaceInteraction {
	panneau := gui.NewPanel(440, 55)
	panneau.SetPosition(740, 760)
	panneau.SetColor4(&math32.Color4{R: .03, G: .03, B: .05, A: .9})

	texte := gui.NewLabel("")
	texte.SetPosition(25, 17)
	panneau.Add(texte)
	panneau.SetVisible(false)
	scene.Add(panneau)
	return &InterfaceInteraction{Panneau: panneau, Texte: texte}
}

func (i *InterfaceInteraction) Afficher(pnj *PNJ) {
	if pnj == nil {
		i.Panneau.SetVisible(false)
		return
	}
	i.Texte.SetText("E — Parler au " + pnj.Nom)
	i.Panneau.SetVisible(true)
}

func NouvelleInterfaceCommerce(scene *core.Node, joueur *library.Character) *InterfaceCommerce {
	panneau := gui.NewPanel(800, 720)
	panneau.SetPosition(560, 170)
	panneau.SetColor4(&math32.Color4{R: .03, G: .03, B: .05, A: .96})

	titre := gui.NewLabel("")
	titre.SetPosition(30, 25)
	titre.SetFontSize(26)
	panneau.Add(titre)

	infos := gui.NewLabel("")
	infos.SetPosition(30, 70)
	panneau.Add(infos)

	message := gui.NewLabel("Choisissez un article.")
	message.SetPosition(30, 610)
	panneau.Add(message)

	boutons := make([]*gui.Button, 0, 10)
	for index := 0; index < 10; index++ {
		bouton := gui.NewButton("")
		bouton.SetSize(740, 45)
		bouton.SetPosition(30, 115+float32(index)*50)
		bouton.SetVisible(false)
		panneau.Add(bouton)
		boutons = append(boutons, bouton)
	}

	fermer := gui.NewButton("FERMER")
	fermer.SetSize(180, 50)
	fermer.SetPosition(590, 650)
	panneau.Add(fermer)

	commerce := &InterfaceCommerce{
		Panneau: panneau, Titre: titre, Infos: infos, Message: message,
		Boutons: boutons, Fermer: fermer, Joueur: joueur,
	}

	for index, bouton := range boutons {
		indexArticle := index
		bouton.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
			var resultat library.ResultatAction
			switch commerce.Mode {
			case PNJMarchand:
				resultat = commerce.Joueur.AcheterMarchand3D(indexArticle)
			case PNJForgeron:
				resultat = commerce.Joueur.FabriquerForgeron3D(indexArticle)
			default:
				return
			}
			commerce.Message.SetText(resultat.Message)
		})
	}

	panneau.SetVisible(false)
	scene.Add(panneau)
	return commerce
}

func (i *InterfaceCommerce) Ouvrir(typePNJ TypePNJ) {
	if i == nil {
		return
	}
	i.Mode = typePNJ
	i.Ouvert = true
	i.Message.SetText("Choisissez un article.")
	if typePNJ == PNJMarchand {
		i.Message.SetText("Les matériaux d'armure se récoltent en combat.")
	}
	i.Panneau.SetVisible(true)
	i.MettreAJour()
}

func (i *InterfaceCommerce) FermerMenu() {
	if i == nil {
		return
	}
	i.Ouvert = false
	i.Panneau.SetVisible(false)
}

func (i *InterfaceCommerce) MettreAJour() {
	if i == nil || !i.Ouvert || i.Joueur == nil {
		return
	}
	i.Infos.SetText(fmt.Sprintf(
		"Or : %d     Inventaire : %d / %d",
		i.Joueur.Argent,
		i.Joueur.TotalInventaire(),
		i.Joueur.CapaciteInventaire,
	))

	switch i.Mode {
	case PNJMarchand:
		i.Titre.SetText("MARCHAND")
		articles := library.BoutiqueMarchand3D()
		for index, bouton := range i.Boutons {
			if index >= len(articles) {
				bouton.SetVisible(false)
				continue
			}
			article := articles[index]
			prix := i.Joueur.PrixPour(article)
			textePrix := fmt.Sprintf("%d pièces d'or", prix)
			if prix == 0 {
				textePrix = "Gratuit"
			}
			bouton.Label.SetText(fmt.Sprintf("%s — %s", article.Nom, textePrix))
			bouton.SetVisible(true)
		}
	case PNJForgeron:
		i.Titre.SetText("FORGERON")
		for index, bouton := range i.Boutons {
			if index >= len(library.Armurerie) {
				bouton.SetVisible(false)
				continue
			}
			equipement := library.Armurerie[index]
			bouton.Label.SetText(fmt.Sprintf(
				"%s — %d or — +%d PV — %s",
				equipement.Nom,
				equipement.Prix,
				equipement.BonusPv,
				texteMateriaux(equipement.Materiaux),
			))
			bouton.SetVisible(true)
		}
	}
}

func texteMateriaux(materiaux map[string]int) string {
	noms := make([]string, 0, len(materiaux))
	for nom := range materiaux {
		noms = append(noms, nom)
	}
	sort.Strings(noms)
	parties := make([]string, 0, len(noms))
	for _, nom := range noms {
		parties = append(parties, fmt.Sprintf("%s x%d", nom, materiaux[nom]))
	}
	return strings.Join(parties, ", ")
}
