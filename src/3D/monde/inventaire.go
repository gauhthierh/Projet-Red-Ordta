package monde

import (
	"fmt"
	"ordta/library"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

const nombreCasesInventaire = 10

// CaseGraphique sépare le fond, l'icône redimensionnée et le texte.
// Cela évite que G3N agrandisse le bouton à la taille native de l'image.
type CaseGraphique struct {
	Panneau *gui.Panel
	Image   *gui.Image
	Texte   *gui.Label
	Nom     string
}

type InterfaceInventaire struct {
	Panneau *gui.Panel
	Ouvert  bool

	Statistiques *gui.Label
	Message      *gui.Label

	Cases          []*CaseGraphique
	ObjetsAffiches []library.ObjetInventaire3D
	CaseChoisie    int

	ArmureTete  *CaseGraphique
	ArmureTorse *CaseGraphique
	ArmurePieds *CaseGraphique

	BoutonUtiliser *gui.Button
	BoutonFermer   *gui.Button
}

func NouvelleInterfaceInventaire(scene *core.Node) *InterfaceInventaire {
	panneau := gui.NewPanel(1800, 930)
	panneau.SetPosition(60, 75)
	panneau.SetColor4(&math32.Color4{R: 0.02, G: 0.02, B: 0.03, A: 0.25})

	fondStatistiques := gui.NewPanel(400, 800)
	fondStatistiques.SetPosition(15, 60)
	fondStatistiques.SetColor4(&math32.Color4{R: 0.03, G: 0.03, B: 0.05, A: 0.90})
	panneau.Add(fondStatistiques)

	fondObjets := gui.NewPanel(990, 800)
	fondObjets.SetPosition(770, 60)
	fondObjets.SetColor4(&math32.Color4{R: 0.03, G: 0.03, B: 0.05, A: 0.90})
	panneau.Add(fondObjets)

	titre := gui.NewLabel("INVENTAIRE")
	titre.SetPosition(800, 25)
	titre.SetFontSize(28)
	panneau.Add(titre)

	statistiques := gui.NewLabel("")
	statistiques.SetPosition(30, 70)
	statistiques.SetFontSize(18)
	panneau.Add(statistiques)

	message := gui.NewLabel("Sélectionnez un objet.")
	message.SetPosition(800, 680)
	panneau.Add(message)

	armureTete := nouvelleCaseGraphique(180, 135, 800, 445)
	armureTorse := nouvelleCaseGraphique(180, 135, 1010, 445)
	armurePieds := nouvelleCaseGraphique(180, 135, 1220, 445)
	panneau.Add(armureTete.Panneau)
	panneau.Add(armureTorse.Panneau)
	panneau.Add(armurePieds.Panneau)

	cases := make([]*CaseGraphique, 0, nombreCasesInventaire)
	for index := 0; index < nombreCasesInventaire; index++ {
		colonne := index % 5
		ligne := index / 5
		caseObjet := nouvelleCaseGraphique(
			150,
			120,
			800+float32(colonne)*175,
			110+float32(ligne)*145,
		)
		panneau.Add(caseObjet.Panneau)
		cases = append(cases, caseObjet)
	}

	boutonUtiliser := gui.NewButton("UTILISER / ÉQUIPER")
	boutonUtiliser.SetSize(260, 55)
	boutonUtiliser.SetPosition(800, 740)
	panneau.Add(boutonUtiliser)

	boutonFermer := gui.NewButton("FERMER [TAB]")
	boutonFermer.SetSize(220, 55)
	boutonFermer.SetPosition(1080, 740)
	panneau.Add(boutonFermer)

	interfaceInventaire := &InterfaceInventaire{
		Panneau:        panneau,
		Statistiques:   statistiques,
		Message:        message,
		Cases:          cases,
		CaseChoisie:    -1,
		ArmureTete:     armureTete,
		ArmureTorse:    armureTorse,
		ArmurePieds:    armurePieds,
		BoutonUtiliser: boutonUtiliser,
		BoutonFermer:   boutonFermer,
	}

	for index, caseObjet := range interfaceInventaire.Cases {
		indexCase := index
		caseObjet.AbonnerClic(func() {
			interfaceInventaire.SelectionnerCase(indexCase)
		})
	}

	panneau.SetVisible(false)
	scene.Add(panneau)
	return interfaceInventaire
}

func nouvelleCaseGraphique(largeur, hauteur, x, y float32) *CaseGraphique {
	fond := gui.NewPanel(largeur, hauteur)
	fond.SetPosition(x, y)
	fond.SetColor4(&math32.Color4{R: 0.18, G: 0.18, B: 0.20, A: 1})

	imageObjet, err := gui.NewImage("../assets/ui/icons/png/inventory.png")
	if err != nil {
		panic(err)
	}
	imageObjet.SetSize(68, 68)
	imageObjet.SetPosition((largeur-68)/2, 8)

	texte := gui.NewLabel("Vide")
	texte.SetPosition(10, hauteur-34)
	texte.SetFontSize(12)

	fond.Add(imageObjet)
	fond.Add(texte)
	return &CaseGraphique{Panneau: fond, Image: imageObjet, Texte: texte}
}

func (c *CaseGraphique) AbonnerClic(action func()) {
	declencher := func(nomEvenement string, evenement interface{}) {
		action()
	}
	c.Panneau.Subscribe(gui.OnMouseDown, declencher)
	c.Image.Subscribe(gui.OnMouseDown, declencher)
	c.Texte.Subscribe(gui.OnMouseDown, declencher)
}

func (c *CaseGraphique) MettreAJour(nom, texte, cheminImage string) {
	c.Texte.SetText(texte)
	if c.Nom == nom {
		return
	}
	if err := c.Image.SetImage(cheminImage); err == nil {
		c.Image.SetSize(68, 68)
	}
	c.Nom = nom
}

func (i *InterfaceInventaire) Ouvrir() {
	if i == nil {
		return
	}
	i.Ouvert = true
	i.CaseChoisie = -1
	i.Panneau.SetVisible(true)
	i.Message.SetText("Sélectionnez un objet.")
}

func (i *InterfaceInventaire) Fermer() {
	if i == nil {
		return
	}
	i.Ouvert = false
	i.CaseChoisie = -1
	i.Panneau.SetVisible(false)
}

func (i *InterfaceInventaire) SelectionnerCase(index int) {
	if i == nil || index < 0 || index >= len(i.ObjetsAffiches) {
		return
	}
	i.CaseChoisie = index
	objet := i.ObjetsAffiches[index]
	i.Message.SetText(objet.Nom)
}

func (i *InterfaceInventaire) ObjetSelectionne() (string, bool) {
	if i == nil || i.CaseChoisie < 0 || i.CaseChoisie >= len(i.ObjetsAffiches) {
		return "", false
	}
	return i.ObjetsAffiches[i.CaseChoisie].Nom, true
}

func (i *InterfaceInventaire) AfficherMessage(message string) {
	if i != nil {
		i.Message.SetText(message)
	}
}

func (i *InterfaceInventaire) MettreAJour(personnage *library.Character) {
	if i == nil || personnage == nil || !i.Ouvert {
		return
	}

	i.Statistiques.SetText(fmt.Sprintf(
		"%s\nClasse : %s\nNiveau : %d\n\nPV : %d / %d\nMana : %d / %d\nAttaque : %d\nInitiative : %d\n\nExpérience : %d / %d\nOr : %d\nInventaire : %d / %d\n\nAttaques : %s\nSorts : %s",
		personnage.Nom,
		personnage.Classe,
		personnage.Niveau,
		personnage.PVActuel,
		personnage.PVMaxTotal,
		personnage.ManaActuel,
		personnage.ManaMax,
		personnage.Attaque,
		personnage.Initiative,
		personnage.ExperienceActuelle,
		personnage.ExperienceMax,
		personnage.Argent,
		personnage.TotalInventaire(),
		personnage.CapaciteInventaire,
		strings.Join(personnage.AttaquesPhysiques, ", "),
		strings.Join(personnage.Skill, ", "),
	))

	i.ObjetsAffiches = personnage.Inventaire3D()
	for index, caseObjet := range i.Cases {
		if index >= len(i.ObjetsAffiches) {
			caseObjet.MettreAJour("", "Vide", "../assets/ui/icons/png/inventory.png")
			continue
		}

		objet := i.ObjetsAffiches[index]
		caseObjet.MettreAJour(
			objet.Nom,
			objet.Nom,
			"../assets/ui/icons/png/"+objet.Icone,
		)
	}

	mettreAJourCaseArmure(i.ArmureTete, "TÊTE", personnage.Equipement.Tete)
	mettreAJourCaseArmure(i.ArmureTorse, "TORSE", personnage.Equipement.Torse)
	mettreAJourCaseArmure(i.ArmurePieds, "PIEDS", personnage.Equipement.Pied)
}

func mettreAJourCaseArmure(caseArmure *CaseGraphique, titre string, equipement library.Stuff) {
	if equipement.Nom == "" {
		caseArmure.MettreAJour("", titre+" — Vide", "../assets/ui/icons/png/inventory.png")
		return
	}
	caseArmure.MettreAJour(
		equipement.Nom,
		fmt.Sprintf("%s — %s  +%d PV", titre, equipement.Nom, equipement.BonusPv),
		"../assets/ui/icons/png/"+library.IconeObjet3D(equipement.Nom),
	)
}
