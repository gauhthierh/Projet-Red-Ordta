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

// TypePNJ : Identifie le service proposé par un personnage non joueur.
type TypePNJ string

const (
	PNJMarchand   TypePNJ = "marchand"
	PNJForgeron   TypePNJ = "forgeron"
	rayonDialogue         = float32(4.5)
)

// PNJ : Associe une position d'interaction, un service et un modèle visuel.
type PNJ struct {
	Type     TypePNJ
	Nom      string
	Position math32.Vector3
	Modele   *personnages.Personnage
}

// InterfaceInteraction : Porte l'indication affichée lorsqu'un PNJ est à portée.
type InterfaceInteraction struct {
	Panneau *gui.Panel
	Texte   *gui.Label
}

// InterfaceCommerce : Relie les boutons de la boutique au même personnage que celui utilisé par le jeu.
type InterfaceCommerce struct {
	Panneau *gui.Panel
	Titre   *gui.Label
	Infos   *gui.Label
	Message *gui.Label
	Boutons []*gui.Button
	Fermer  *gui.Button

	Page       int
	Pagination *gui.Label
	Ouvert     bool
	Mode       TypePNJ
	Joueur     *library.Character
}

// AjouterPNJMarche : Installe le marchand et le forgeron avec leur tenue propre aux coordonnées de leurs cabanons.
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

// PNJLePlusProche : Cherche le PNJ le plus proche dans le rayon d'interaction, en ignorant la hauteur.
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

// NouvelleInterfaceInteraction : Crée l'indication de proximité invitant à parler avec E.
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

// Afficher : Affiche le nom du PNJ disponible, ou cache l'indication si aucun n'est à portée.
func (i *InterfaceInteraction) Afficher(pnj *PNJ) {
	if pnj == nil {
		i.Panneau.SetVisible(false)
		return
	}
	i.Texte.SetText("E — Parler au " + pnj.Nom)
	i.Panneau.SetVisible(true)
}

// NouvelleInterfaceCommerce : Construit la boutique et branche les achats ou fabrications sur library_3d.
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

	boutons := make([]*gui.Button, 0, 8)
	for index := 0; index < 8; index++ {
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
				resultat = commerce.Joueur.AcheterMarchand3D(commerce.Page*8 + indexArticle)
			case PNJForgeron:
				resultat = commerce.Joueur.FabriquerForgeron3D(indexArticle)
			default:
				return
			}
			commerce.Message.SetText(resultat.Message)
		})
	}

	commerce.Pagination = gui.NewLabel("")
	commerce.Pagination.SetPosition(330, 550)
	panneau.Add(commerce.Pagination)
	for index, sens := range []int{-1, 1} {
		direction := sens
		texte := "PRÉCÉDENT"
		if sens > 0 {
			texte = "SUIVANT"
		}
		bouton := gui.NewButton(texte)
		bouton.SetSize(180, 40)
		bouton.SetPosition(30+float32(index)*560, 540)
		bouton.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
			pages := 1
			if commerce.Mode == PNJMarchand {
				pages = (len(library.BoutiqueMarchand3D()) + 7) / 8
			}
			page := commerce.Page + direction
			if page >= 0 && page < pages {
				commerce.Page = page
				commerce.MettreAJour()
			}
		})
		panneau.Add(bouton)
	}
	panneau.SetVisible(false)
	scene.Add(panneau)
	return commerce
}

// Ouvrir : Choisit la boutique du PNJ, revient à la première page et actualise son contenu.
func (i *InterfaceCommerce) Ouvrir(typePNJ TypePNJ) {
	if i == nil {
		return
	}
	i.Mode = typePNJ
	i.Page = 0
	i.Ouvert = true
	i.Message.SetText("Choisissez un article.")
	if typePNJ == PNJMarchand {
		i.Message.SetText("Les achats respectent le niveau indiqué.")
	}
	i.Panneau.SetVisible(true)
	i.MettreAJour()
}

// FermerMenu : Masque la boutique sans changer les objets ni l'or du joueur.
func (i *InterfaceCommerce) FermerMenu() {
	if i == nil {
		return
	}
	i.Ouvert = false
	i.Panneau.SetVisible(false)
}

// MettreAJour : Relit les prix, les recettes et l'état du joueur sans effectuer d'achat.
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

	pages := 1
	if i.Mode == PNJMarchand {
		pages = (len(library.BoutiqueMarchand3D()) + 7) / 8
	}
	i.Pagination.SetText(fmt.Sprintf("Page %d / %d", i.Page+1, pages))
	switch i.Mode {
	case PNJMarchand:
		i.Titre.SetText("MARCHAND")
		articles := library.BoutiqueMarchand3D()
		for numero, bouton := range i.Boutons {
			index := i.Page*8 + numero
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
			bouton.Label.SetText(fmt.Sprintf("%s — %s — niv. %d", article.Nom, textePrix, article.NiveauMin))
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

// texteMateriaux : Trie les matériaux par nom afin que les recettes gardent un affichage stable.
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
