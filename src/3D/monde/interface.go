package monde

import (
	"fmt"
	"ordta/library"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type InterfaceCombat struct {
	Panneau          *gui.Panel
	BoutonAttaquer   *gui.Button
	BoutonSorts      *gui.Button
	BoutonInventaire *gui.Button
	BoutonDefendre   *gui.Button
	BoutonQuitter    *gui.Button
	BoutonRejouer    *gui.Button
	BoutonsChoix     []*gui.Button
	ModeChoix        string
	OptionsChoix     []string
}

type InterfacePersonnage struct {
	Panneau *gui.Panel

	FondVie        *gui.Panel
	RemplissageVie *gui.Panel
	TexteVie       *gui.Label
	LargeurVieMax  float32

	NomClasse  *gui.Label
	Mana       *gui.Label
	Combat     *gui.Label
	Niveau     *gui.Label
	Experience *gui.Label
	Argent     *gui.Label
}

type InterfaceMonstre struct {
	Panneau *gui.Panel

	Nom          *gui.Label
	Statistiques *gui.Label
	BoutonCibler *gui.Button

	FondVie        *gui.Panel
	RemplissageVie *gui.Panel
	TexteVie       *gui.Label
	LargeurVieMax  float32
}

type InterfaceMonstres struct {
	Panneau  *gui.Panel
	Monstres []*InterfaceMonstre
}

type InterfaceEtatCombat struct {
	Panneau *gui.Panel
	Vague   *gui.Label
	Tour    *gui.Label
	Phase   *gui.Label
	Message *gui.Label
	Butin   *gui.Label
}

func NouvelleInterfaceCombat(scene *core.Node) *InterfaceCombat {
	panneau := gui.NewPanel(800, 230)
	panneau.SetPosition(560, 830)
	panneau.SetColor4(&math32.Color4{R: 0.03, G: 0.03, B: 0.05, A: 0.90})

	bouton := gui.NewButton("ATTAQUES")
	bouton.SetSize(175, 50)
	bouton.SetPosition(15, 165)

	boutonSorts := gui.NewButton("SORTS")
	boutonSorts.SetSize(175, 50)
	boutonSorts.SetPosition(210, 165)

	boutonInventaire := gui.NewButton("INVENTAIRE")
	boutonInventaire.SetSize(175, 50)
	boutonInventaire.SetPosition(405, 165)

	boutonDefendre := gui.NewButton("DÉFENDRE")
	boutonDefendre.SetSize(175, 50)
	boutonDefendre.SetPosition(600, 165)

	boutonQuitter := gui.NewButton("QUITTER L'ARÈNE")
	boutonQuitter.SetSize(220, 55)
	boutonQuitter.SetPosition(565, 110)
	boutonQuitter.SetVisible(false)
	boutonRejouer := gui.NewButton("CHOISIR UN AUTRE COMBAT")
	boutonRejouer.SetPosition(15, 110)
	boutonRejouer.SetSize(260, 50)
	boutonRejouer.SetVisible(false)
	panneau.Add(boutonRejouer)

	panneau.Add(bouton)
	panneau.Add(boutonSorts)
	panneau.Add(boutonInventaire)
	panneau.Add(boutonDefendre)
	panneau.Add(boutonQuitter)

	boutonsChoix := make([]*gui.Button, 0, 12)
	for index := 0; index < 12; index++ {
		colonne := index % 6
		ligne := index / 6
		boutonChoix := gui.NewButton("")
		boutonChoix.Label.SetFontSize(11)
		boutonChoix.SetSize(120, 45)
		boutonChoix.SetPosition(15+float32(colonne)*130, 10+float32(ligne)*55)
		boutonChoix.SetVisible(false)
		panneau.Add(boutonChoix)
		boutonsChoix = append(boutonsChoix, boutonChoix)
	}

	panneau.SetVisible(false)

	scene.Add(panneau)

	return &InterfaceCombat{
		Panneau:          panneau,
		BoutonAttaquer:   bouton,
		BoutonSorts:      boutonSorts,
		BoutonInventaire: boutonInventaire,
		BoutonDefendre:   boutonDefendre,
		BoutonQuitter:    boutonQuitter,
		BoutonRejouer:    boutonRejouer,
		BoutonsChoix:     boutonsChoix,
	}
}

func (i *InterfaceCombat) AfficherChoix(mode string, options []string) {
	if i == nil {
		return
	}
	i.ModeChoix = mode
	i.OptionsChoix = append([]string(nil), options...)
	for index, bouton := range i.BoutonsChoix {
		if index >= len(options) {
			bouton.SetVisible(false)
			continue
		}
		bouton.Label.SetText(options[index])
		if mode == "sort" {
			bouton.Label.SetText(fmt.Sprintf("%s\n%d mana", options[index], library.CoutSort3D(options[index])))
		}
		bouton.SetVisible(true)
	}
}

func (i *InterfaceCombat) MasquerChoix() {
	if i == nil {
		return
	}
	i.ModeChoix = ""
	i.OptionsChoix = nil
	for _, bouton := range i.BoutonsChoix {
		bouton.SetVisible(false)
	}
}

func NouvelleInterfacePersonnage(scene *core.Node) *InterfacePersonnage {
	// Panneau principal de l'interface
	panneau := gui.NewPanel(380, 230)
	panneau.SetPosition(30, 30)
	panneau.SetColor4(&math32.Color4{
		R: 0.05,
		G: 0.05,
		B: 0.08,
		A: 0.90,
	})

	// Nom et classe du personnage
	nomClasse := gui.NewLabel("")
	nomClasse.SetPosition(15, 12)

	// Fond de la barre de vie
	fondVie := gui.NewPanel(340, 22)
	fondVie.SetPosition(20, 42)
	fondVie.SetColor4(&math32.Color4{
		R: 0.20,
		G: 0.03,
		B: 0.03,
		A: 1,
	})

	// Partie rouge représentant les points de vie restants
	remplissageVie := gui.NewPanel(340, 22)
	remplissageVie.SetPosition(20, 42)
	remplissageVie.SetColor4(&math32.Color4{
		R: 0.80,
		G: 0.08,
		B: 0.08,
		A: 1,
	})

	// Texte affiché par-dessus la barre de vie
	texteVie := gui.NewLabel("")
	texteVie.SetPosition(145, 44)

	// Mana
	mana := gui.NewLabel("")
	mana.SetPosition(15, 78)

	// Attaque et initiative
	combat := gui.NewLabel("")
	combat.SetPosition(15, 102)

	// Niveau
	niveau := gui.NewLabel("")
	niveau.SetPosition(15, 126)

	// Expérience
	experience := gui.NewLabel("")
	experience.SetPosition(15, 150)

	// Argent et remplissage de l'inventaire
	argent := gui.NewLabel("")
	argent.SetPosition(15, 174)

	// L'ordre d'ajout est important :
	// le fond doit être derrière le remplissage et le texte.
	panneau.Add(fondVie)
	panneau.Add(remplissageVie)
	panneau.Add(texteVie)
	panneau.Add(nomClasse)
	panneau.Add(mana)
	panneau.Add(combat)
	panneau.Add(niveau)
	panneau.Add(experience)
	panneau.Add(argent)

	// On ajoute uniquement le panneau principal à la scène.
	scene.Add(panneau)

	return &InterfacePersonnage{
		Panneau: panneau,

		FondVie:        fondVie,
		RemplissageVie: remplissageVie,
		TexteVie:       texteVie,
		LargeurVieMax:  340,

		NomClasse:  nomClasse,
		Mana:       mana,
		Combat:     combat,
		Niveau:     niveau,
		Experience: experience,
		Argent:     argent,
	}
}

func (i *InterfacePersonnage) MettreAJourPersonnage(p *library.Character) {
	if i == nil || p == nil {
		return
	}

	// Calcul de la proportion de vie restante.
	proportionVie := float32(0)

	if p.PvMaxTotal > 0 {
		proportionVie = float32(p.PvActuel) / float32(p.PvMaxTotal)
	}

	// Empêche la barre de dépasser ses limites.
	if proportionVie < 0 {
		proportionVie = 0
	}

	if proportionVie > 1 {
		proportionVie = 1
	}

	// Modification de la largeur de la barre rouge.
	largeurVie := i.LargeurVieMax * proportionVie
	i.RemplissageVie.SetSize(largeurVie, 22)

	// Mise à jour des textes.
	i.TexteVie.SetText(
		fmt.Sprintf("PV : %d / %d", p.PvActuel, p.PvMaxTotal),
	)

	i.NomClasse.SetText(
		fmt.Sprintf("%s - %s", p.Nom, p.Classe),
	)

	i.Mana.SetText(
		fmt.Sprintf("Mana : %d / %d", p.ManaActuel, p.ManaMax),
	)

	i.Combat.SetText(
		fmt.Sprintf("Attaque : %d | Initiative : %d", p.Attaque, p.Initiative),
	)

	i.Niveau.SetText(
		fmt.Sprintf("Niveau : %d", p.Niveau),
	)

	i.Experience.SetText(
		fmt.Sprintf("Expérience : %d / %d", p.ExperienceActuelle, p.ExperienceMax),
	)

	i.Argent.SetText(
		fmt.Sprintf("Or : %d", p.Argent),
	)
}

func NouvelleFicheMonstre(positionY float32) *InterfaceMonstre {
	// Conteneur de la fiche
	panneau := gui.NewPanel(320, 100)
	panneau.SetPosition(0, positionY)
	panneau.SetColor4(&math32.Color4{
		R: 0.05,
		G: 0.05,
		B: 0.08,
		A: 0.90,
	})

	// Nom
	nom := gui.NewLabel("Nom du monstre")
	nom.SetPosition(15, 10)

	// Fond de la barre de vie
	fondVie := gui.NewPanel(280, 20)
	fondVie.SetPosition(20, 36)
	fondVie.SetColor4(&math32.Color4{
		R: 0.20,
		G: 0.03,
		B: 0.03,
		A: 1,
	})

	// Vie restante
	remplissageVie := gui.NewPanel(280, 20)
	remplissageVie.SetPosition(20, 36)
	remplissageVie.SetColor4(&math32.Color4{
		R: 0.75,
		G: 0.08,
		B: 0.08,
		A: 1,
	})

	// Valeur des PV affichée sur la barre
	texteVie := gui.NewLabel("")
	texteVie.SetPosition(110, 37)

	// Autres statistiques
	statistiques := gui.NewLabel("")
	statistiques.SetPosition(15, 68)

	boutonCibler := gui.NewButton("CIBLER")
	boutonCibler.SetSize(85, 28)
	boutonCibler.SetPosition(215, 64)

	// Ajout dans l’ordre d’affichage
	panneau.Add(fondVie)
	panneau.Add(remplissageVie)
	panneau.Add(texteVie)
	panneau.Add(nom)
	panneau.Add(statistiques)
	panneau.Add(boutonCibler)

	return &InterfaceMonstre{
		Panneau:        panneau,
		Nom:            nom,
		Statistiques:   statistiques,
		BoutonCibler:   boutonCibler,
		FondVie:        fondVie,
		RemplissageVie: remplissageVie,
		TexteVie:       texteVie,
		LargeurVieMax:  280,
	}
}

func NouvelleInterfaceMonstres(scene *core.Node, nombreMaximum int) *InterfaceMonstres {
	// Conteneur général placé à droite de l'écran.
	hauteur := float32(nombreMaximum) * 115

	panneau := gui.NewPanel(320, hauteur)
	panneau.SetPosition(1570, 30)

	// Tableau qui conservera les références vers chaque fiche.
	fiches := make([]*InterfaceMonstre, 0, nombreMaximum)

	for index := 0; index < nombreMaximum; index++ {
		// Chaque fiche est décalée de 115 pixels vers le bas.
		positionY := float32(index) * 115

		fiche := NouvelleFicheMonstre(positionY)

		panneau.Add(fiche.Panneau)
		fiches = append(fiches, fiche)
	}

	// Les monstres ne doivent pas être affichés en exploration.
	panneau.SetVisible(false)

	scene.Add(panneau)

	return &InterfaceMonstres{
		Panneau:  panneau,
		Monstres: fiches,
	}
}

func (i *InterfaceMonstres) MettreAJourMonstres(ennemis []library.EnnemiCombat) {
	if i == nil {
		return
	}

	// Le panneau général est visible uniquement s'il existe des ennemis.
	i.Panneau.SetVisible(len(ennemis) > 0)

	for index, fiche := range i.Monstres {
		// Il peut y avoir plus de fiches préparées que de monstres.
		if index >= len(ennemis) {
			fiche.Panneau.SetVisible(false)
			continue
		}

		fiche.Panneau.SetVisible(true)

		monstre := &ennemis[index].Monstre

		// Calcul de la proportion de vie restante.
		proportionVie := float32(0)

		if monstre.PvMax > 0 {
			proportionVie =
				float32(monstre.PvActuel) /
					float32(monstre.PvMax)
		}

		if proportionVie < 0 {
			proportionVie = 0
		}

		if proportionVie > 1 {
			proportionVie = 1
		}

		// Mise à jour de la barre.
		largeurVie := fiche.LargeurVieMax * proportionVie
		fiche.RemplissageVie.SetSize(largeurVie, 20)

		// Mise à jour des informations.
		fiche.Nom.SetText(monstre.Nom)

		fiche.TexteVie.SetText(
			fmt.Sprintf("PV : %d / %d", monstre.PvActuel, monstre.PvMax),
		)

		fiche.Statistiques.SetText(
			fmt.Sprintf("Attaque : %d | Initiative : %d", monstre.Attaque, monstre.Initiative),
		)
	}
}

func NouvelleInterfaceEtatCombat(scene *core.Node) *InterfaceEtatCombat {
	panneau := gui.NewPanel(780, 235)
	panneau.SetPosition(660, 30)
	panneau.SetColor4(&math32.Color4{
		R: 0.05,
		G: 0.05,
		B: 0.08,
		A: 0.90,
	})

	vague := gui.NewLabel("")
	vague.SetPosition(20, 12)

	tour := gui.NewLabel("")
	tour.SetPosition(200, 12)

	phase := gui.NewLabel("")
	phase.SetPosition(350, 12)

	message := gui.NewLabel("")
	message.SetPosition(20, 65)

	panneau.Add(vague)
	panneau.Add(tour)
	panneau.Add(phase)
	panneau.Add(message)
	butin := gui.NewLabel("")
	butin.SetPosition(20, 115)
	panneau.Add(butin)

	panneau.SetVisible(false)
	scene.Add(panneau)

	return &InterfaceEtatCombat{
		Panneau: panneau,
		Vague:   vague,
		Tour:    tour,
		Phase:   phase,
		Message: message,
		Butin:   butin,
	}
}

func textePhaseCombat(phase library.PhaseCombat) string {
	switch phase {
	case library.PhaseTourJoueur:
		return "Tour du joueur"

	case library.PhaseTourMonstres:
		return "Tour des monstres"

	case library.PhaseEntreVagues:
		return "Vague terminée"

	case library.PhaseVictoire:
		return "Victoire"

	case library.PhaseDefaite:
		return "Défaite"

	default:
		return "Phase inconnue"
	}
}

func (i *InterfaceEtatCombat) MettreAJour(combat *library.CombatArene, dernierMessage string) {
	if i == nil {
		return
	}

	if combat == nil {
		i.Panneau.SetVisible(false)
		return
	}

	i.Panneau.SetVisible(true)

	i.Vague.SetText(
		fmt.Sprintf("%s : %d / %d", combat.Mode, combat.NumeroVague, len(combat.Vagues)),
	)

	i.Tour.SetText(
		fmt.Sprintf("Tour : %d", combat.Tour),
	)

	i.Phase.SetText(
		textePhaseCombat(combat.Phase),
	)

	if dernierMessage == "" {
		dernierMessage = "Préparez-vous au combat."
	}

	// Couper les messages longs pour qu'un butin ne déborde pas du panneau.
	texte, ligne := "", ""
	for _, mot := range strings.Fields(dernierMessage) {
		if len([]rune(ligne+mot)) > 85 {
			texte += strings.TrimSpace(ligne) + "\n"
			ligne = ""
		}
		ligne += mot + " "
	}
	i.Message.SetText(texte + strings.TrimSpace(ligne))
	i.Butin.SetText(fmt.Sprintf("Or gagné : %d\nButin :\n", combat.OrTotal) + strings.Join(combat.Butins, "\n"))
}
