// Package monde contient la scène 3D et sa boucle de jeu.
package monde

import (
	"fmt"
	"time"

	"ordta/3D/monstres"
	"ordta/3D/personnages"
	"ordta/library"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

const (
	vitesseMarche   float32 = 10 // Unités parcourues par seconde.
	vitesseSprint   float32 = 25
	rayonPersonnage float32 = 0.3
)

// Lancer construit le monde, puis démarre sa boucle d'affichage.
func Lancer() {
	// Fenêtre G3N et racine de la scène 3D.
	ordta := app.App()
	sons := NouvelAudioJeu("../assets/audio")
	// Libérer les sons avant que G3N ne ferme le périphérique audio.
	ordta.Subscribe(app.OnExit, func(_ string, _ interface{}) { sons.Fermer() })
	scene := core.NewNode()

	// Gestion de la map
	CheminOBJ := "../assets/maps/red_world/red_world_map_3d.obj"
	CheminMTL := "../assets/maps/red_world/red_world_map_3d.mtl"

	donneesMonde, err := ChargerZones("../assets/maps/red_world/red_world_layout.json")
	if err != nil {
		panic(err)
	}

	spawn := donneesMonde.Spawn

	maps3d, err := ChargerMap3D(CheminOBJ, CheminMTL)

	if err != nil {
		panic(err)
	}
	scene.Add(maps3d)
	AjouterCabanons(scene)

	// Marchand et forgeron placés devant les deux premiers cabanons du marché.
	pnjs := AjouterPNJMarche(scene)
	for index, pnj := range pnjs {
		donneesMonde.Obstacles = append(donneesMonde.Obstacles, Collision{
			Identifiant: fmt.Sprintf("pnj_%d", index),
			Categorie:   "pnj",
			Forme:       "circle",
			X:           pnj.Position.X,
			Y:           pnj.Position.Y,
			Rayon:       .65,
		})
	}

	// Arène
	combatEnCours := false

	// Personnage
	PersonnageBackend := library.NouveauPersonnage3D("Joueur", "Humain")
	partieChargee, erreurSauvegarde := ChargerPartie3D(cheminSauvegarde)
	if partieChargee != nil {
		PersonnageBackend = partieChargee.Personnage
		position := partieChargee.Position
		if positionAutorisee(position.X, position.Y, rayonPersonnage, donneesMonde) && !EstDansArene(position.X, position.Y) {
			spawn = [3]float32{position.X, position.Y, position.Z}
		}
	}

	var combatBackend *library.CombatArene
	effetsPersonnage := library.EffetsPersonnage3D{Joueur: &PersonnageBackend}
	cibleSelectionnee := 0
	tabEtaitAppuye := false
	eEtaitAppuye := false
	var pnjActif *PNJ

	// Chargement des monstres
	var modelesAdversaires []*monstres.Monstre

	// Chronomètre
	tempsAvantActionMonstre := time.Duration(0)
	tempsAvantVagueSuivante := time.Duration(0)

	// GUI / Interface
	interfaceJeu := core.NewNode()
	scene.Add(interfaceJeu)
	interfaceCombat := NouvelleInterfaceCombat(interfaceJeu)
	interfacePersonnage := NouvelleInterfacePersonnage(interfaceJeu)
	interfaceMonstres := NouvelleInterfaceMonstres(interfaceJeu, len(PositionsMonstresArene))
	interfaceEtatCombat := NouvelleInterfaceEtatCombat(interfaceJeu)
	journalCombat := NouveauJournalCombat(interfaceJeu)
	animationCombat := NouvelleAnimationCombat(scene)
	interfaceInventaire := NouvelleInterfaceInventaire(interfaceJeu)
	interfaceInteraction := NouvelleInterfaceInteraction(interfaceJeu)
	interfaceCommerce := NouvelleInterfaceCommerce(interfaceJeu, &PersonnageBackend)
	choixCombat := NouveauChoixCombat(interfaceJeu)
	dernierMessageCombat := ""

	gui.Manager().Set(scene)

	// Le package personnages fournit le modèle détaillé.
	// Son nœud de déplacement reste Z-up, comme le sol et la caméra.
	personnage3d := personnages.Nouveau()
	noeudPersonnage := personnage3d.Noeud()
	scene.Add(noeudPersonnage)
	noeudPersonnage.SetPosition(spawn[0], spawn[1], spawn[2])

	// Caméra
	cameraSimulation := camera.New(1920.0 / 1080.0)

	ConfigurerFenetreSimulation(
		ordta,
		cameraSimulation,
	)

	largeurMenu, hauteurMenu := ordta.GetSize()
	menu := NouveauMenuJeu(scene, float32(largeurMenu), float32(hauteurMenu))
	creation := NouvelleCreationPersonnage(scene, float32(largeurMenu), float32(hauteurMenu))
	if partieChargee != nil {
		menu.BoutonDemarrer.Label.SetText("CONTINUER")
	}
	if erreurSauvegarde != nil {
		menu.BoutonDemarrer.SetEnabled(false)
		message := gui.NewLabel("Sauvegarde illisible : elle est conservée. Consultez le terminal.")
		message.SetPosition(30, 40)
		menu.Panneau.Add(message)
		fmt.Println(erreurSauvegarde)
	}
	interfaceJeu.SetVisible(false)
	gui.Manager().Set(menu.Panneau)
	LibererSourisSimulation()
	echapEtaitAppuye := false
	angleHorizontal := float32(0)
	angleVertical := float32(0)
	if partieChargee != nil {
		angleHorizontal = partieChargee.AngleHorizontal
		angleVertical = partieChargee.AngleVertical
	}
	tempsSauvegarde := time.Duration(0)
	f5EtaitAppuye := false
	sauvegarder := func() {
		if combatEnCours || !menu.Demarre || effetsPersonnage.Actif() {
			return
		}
		err := SauvegarderPartie3D(cheminSauvegarde, Partie3D{
			Personnage: PersonnageBackend, Position: noeudPersonnage.Position(),
			AngleHorizontal: angleHorizontal, AngleVertical: angleVertical,
		})
		if err != nil {
			fmt.Println("Échec sauvegarde :", err)
		}
	}
	defer func() {
		effetsPersonnage.MettreAJour(3 * time.Second)
		sauvegarder()
	}()

	// Reprendre conserve l'inventaire ou le commerce qui était ouvert.
	reprendreJeu := func() {
		menu.Reprendre()
		interfaceJeu.SetVisible(true)
		gui.Manager().Set(interfaceJeu)
		if !combatEnCours && !interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert {
			VerrouillerSourisSimulation()
		} else {
			LibererSourisSimulation()
		}
	}
	menu.BoutonReprendre.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		reprendreJeu()
	})
	menu.BoutonDemarrer.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		if erreurSauvegarde != nil {
			return
		}
		if partieChargee != nil {
			reprendreJeu()
			return
		}
		menu.Panneau.SetVisible(false)
		creation.Panneau.SetVisible(true)
		gui.Manager().Set(creation.Panneau)
		gui.Manager().SetKeyFocus(creation.Nom)
	})
	creation.Retour.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		creation.Panneau.SetVisible(false)
		menu.Panneau.SetVisible(true)
		gui.Manager().SetKeyFocus(nil)
		gui.Manager().Set(menu.Panneau)
	})
	creation.Valider.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		personnage, err := library.CreerPersonnage3D(creation.Nom.Text(), creation.Classe)
		if err != nil {
			creation.Message.SetText(err.Error())
			return
		}
		PersonnageBackend = personnage
		creation.Panneau.SetVisible(false)
		gui.Manager().SetKeyFocus(nil)
		reprendreJeu()
		sauvegarder()
	})
	menu.BoutonQuitter.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		ordta.Exit() // La sauvegarde normale est effectuée par defer.
	})
	PlacerCameraPersonnage(cameraSimulation, noeudPersonnage, angleHorizontal, angleVertical)

	ActiverRegardSouris(ordta, &angleHorizontal, &angleVertical, func() bool {
		return !menu.Ouvert && !combatEnCours && !interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert
	})

	fermerInventaire := func() {
		interfaceInventaire.Fermer()
		sauvegarder()
		if !combatEnCours {
			VerrouillerSourisSimulation()
		}
	}

	ouvrirInventaire := func() {
		if animationCombat.Active {
			return
		}
		if interfaceCommerce.Ouvert {
			return
		}
		if combatEnCours &&
			(combatBackend == nil || combatBackend.Phase != library.PhaseTourJoueur) {
			return
		}
		interfaceCombat.MasquerChoix()
		interfaceInventaire.Ouvrir()
		LibererSourisSimulation()
	}

	fermerCommerce := func() {
		interfaceCommerce.FermerMenu()
		sauvegarder()
		pnjActif = nil
		VerrouillerSourisSimulation()
	}

	appliquerResultat := func(action string, resultat library.ResultatAction, indexCible int) {
		dernierMessageCombat = resultat.Message
		journalCombat.Ajouter(texteResultatCombat(action, resultat))
		if resultat.Reussite && resultat.TourConsomme {
			cible := noeudPersonnage
			if indexCible >= 0 && indexCible < len(modelesAdversaires) && resultat.Degats > 0 {
				cible = modelesAdversaires[indexCible].Noeud()
			}
			animationCombat.Demarrer(action, noeudPersonnage, cible, personnage3d.MembresCombat(), func() {
				if resultat.CibleVaincue && indexCible >= 0 && indexCible < len(modelesAdversaires) {
					modelesAdversaires[indexCible].Noeud().SetVisible(false)
				}
			})
			interfaceCombat.MasquerChoix()
		}
	}

	// Sélection d'une cible en cliquant sur sa fiche à droite.
	for index, fiche := range interfaceMonstres.Monstres {
		indexMonstre := index
		fiche.BoutonCibler.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
			if combatBackend == nil || indexMonstre >= len(combatBackend.Ennemis) {
				return
			}
			if combatBackend.Ennemis[indexMonstre].Monstre.PvActuel <= 0 {
				return
			}
			cibleSelectionnee = indexMonstre
			dernierMessageCombat = "Cible : " + combatBackend.Ennemis[indexMonstre].Monstre.Nom
		})
	}

	interfaceCombat.BoutonAttaquer.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if !animationCombat.Active && combatBackend != nil && combatBackend.Phase == library.PhaseTourJoueur {
			interfaceCombat.AfficherChoix("attaque", combatBackend.AttaquesDisponibles())
		}
	})

	interfaceCombat.BoutonSorts.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if !animationCombat.Active && combatBackend != nil && combatBackend.Phase == library.PhaseTourJoueur {
			interfaceCombat.AfficherChoix("sort", combatBackend.SortsDisponibles())
		}
	})

	interfaceCombat.BoutonInventaire.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		ouvrirInventaire()
	})

	interfaceCombat.BoutonDefendre.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if combatBackend == nil || animationCombat.Active {
			return
		}
		resultat := combatBackend.Defendre()
		appliquerResultat("defense", resultat, -1)
	})

	for index, bouton := range interfaceCombat.BoutonsChoix {
		indexOption := index
		bouton.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
			if animationCombat.Active || combatBackend == nil ||
				combatBackend.Phase != library.PhaseTourJoueur ||
				indexOption >= len(interfaceCombat.OptionsChoix) {
				return
			}

			cibles := combatBackend.CiblesVivantes()
			if len(cibles) == 0 {
				return
			}
			if cibleSelectionnee >= len(combatBackend.Ennemis) ||
				combatBackend.Ennemis[cibleSelectionnee].Monstre.PvActuel <= 0 {
				cibleSelectionnee = cibles[0]
			}

			option := interfaceCombat.OptionsChoix[indexOption]
			var resultat library.ResultatAction
			switch interfaceCombat.ModeChoix {
			case "attaque":
				resultat = combatBackend.AttaquerPhysiquement(option, cibleSelectionnee)
			case "sort":
				resultat = combatBackend.LancerSort(option, cibleSelectionnee)
			default:
				return
			}
			if resultat.Reussite {
				if interfaceCombat.ModeChoix == "sort" {
					sons.Effet("sort")
				} else {
					sons.Effet("attaque")
				}
			}
			appliquerResultat(option, resultat, cibleSelectionnee)
		})
	}

	interfaceInventaire.BoutonFermer.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		fermerInventaire()
	})

	interfaceCommerce.Fermer.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		fermerCommerce()
	})

	interfaceInventaire.BoutonUtiliser.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		nomObjet, existe := interfaceInventaire.ObjetSelectionne()
		if !existe {
			interfaceInventaire.AfficherMessage("Sélectionnez d'abord un objet.")
			return
		}

		if combatEnCours {
			if combatBackend == nil || animationCombat.Active {
				return
			}
			resultat := combatBackend.UtiliserObjet(nomObjet)
			if resultat.Reussite {
				sons.Effet("potion")
			}
			interfaceInventaire.AfficherMessage(resultat.Message)
			journalCombat.Ajouter("Objet utilisé : " + nomObjet)
			appliquerResultat("objet", resultat, -1)
			if resultat.Reussite && resultat.TourConsomme {
				fermerInventaire()
			}
			return
		}

		resultat := effetsPersonnage.Utiliser(nomObjet)
		if resultat.Reussite {
			sons.Effet("potion")
		}
		interfaceInventaire.AfficherMessage(resultat.Message)
	})

	// Cliquer sur un emplacement équipé retire son objet.
	for emplacement, caseArmure := range map[string]*CaseGraphique{
		library.EmplacementTete:  interfaceInventaire.ArmureTete,
		library.EmplacementTorse: interfaceInventaire.ArmureTorse,
		library.EmplacementPied:  interfaceInventaire.ArmurePieds,
	} {
		nomEmplacement := emplacement
		caseArmure.AbonnerClic(func() {
			if combatEnCours {
				interfaceInventaire.AfficherMessage("L'équipement ne peut pas être changé pendant le combat.")
				return
			}
			resultat := PersonnageBackend.Desequiper3D(nomEmplacement)
			interfaceInventaire.AfficherMessage(resultat.Message)
		})
	}

	// Une seule remise à zéro, partagée par la sortie et le choix du combat suivant.
	nettoyerCombat := func() {
		animationCombat.Arreter()
		journalCombat.Reinitialiser()
		if combatBackend != nil {
			combatBackend.Quitter3D()
		}
		for _, modele := range modelesAdversaires {
			scene.Remove(modele.Noeud())
		}
		modelesAdversaires = nil
		combatBackend = nil
		cibleSelectionnee = 0
		dernierMessageCombat = ""
		tempsAvantActionMonstre = 0
		tempsAvantVagueSuivante = 0
		interfaceCombat.MasquerChoix()
		interfaceInventaire.Fermer()
		sons.ArreterEffets()
	}
	quitterArene := func() {
		nettoyerCombat()
		combatEnCours = false
		choixCombat.Panneau.SetVisible(false)
		noeudPersonnage.SetPosition(PositionSortieArene.X, PositionSortieArene.Y, PositionSortieArene.Z)
		sauvegarder()
		VerrouillerSourisSimulation()
	}
	interfaceCombat.BoutonQuitter.Subscribe(gui.OnClick, func(_ string, _ interface{}) { quitterArene() })
	choixCombat.Quitter.Subscribe(gui.OnClick, func(_ string, _ interface{}) { quitterArene() })
	interfaceCombat.BoutonRejouer.Subscribe(gui.OnClick, func(_ string, _ interface{}) {
		if combatBackend == nil || (combatBackend.Phase != library.PhaseVictoire && combatBackend.Phase != library.PhaseDefaite) {
			return
		}
		nettoyerCombat()
		choixCombat.Ouvrir()
	})
	demarrerCombat := func(mode library.ModeCombat, genre library.TypeMonstre) {
		if effetsPersonnage.Actif() {
			choixCombat.Message.SetText("Attendez la fin du poison.")
			return
		}
		if combatBackend != nil {
			return
		}
		nouveau, err := library.NouveauCombat3D(&PersonnageBackend, mode, genre)
		if err != nil {
			choixCombat.Message.SetText(err.Error())
			return
		}
		modeles, err := ChargerMontres(scene, nouveau.Ennemis)
		if err != nil {
			nouveau.Quitter3D()
			choixCombat.Message.SetText(err.Error())
			return
		}
		combatBackend = nouveau
		journalCombat.Reinitialiser()
		journalCombat.Etat(nouveau)
		modelesAdversaires = modeles
		choixCombat.Panneau.SetVisible(false)
		dernierMessageCombat = "L’initiative détermine le premier tour. Vous pouvez quitter à tout moment."
	}
	choixCombat.Entrainement.Subscribe(gui.OnClick, func(_ string, _ interface{}) { demarrerCombat(library.ModeEntrainement, library.TypeGobelin) })
	choixCombat.Arene.Subscribe(gui.OnClick, func(_ string, _ interface{}) { demarrerCombat(library.ModeArene, library.TypeGobelin) })
	for genre, bouton := range choixCombat.Creatures {
		typeChoisi := genre
		bouton.Subscribe(gui.OnClick, func(_ string, _ interface{}) { demarrerCombat(library.ModeDuel, typeChoisi) })
	}

	// Lumière dans le jeu
	scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 1))

	// Boucle principale : mise à jour du jeu, puis affichage de chaque image.
	ordta.Gls().ClearColor(0.15, 0.25, 0.30, 1)
	ordta.Run(func(rendu *renderer.Renderer, tempsImage time.Duration) {
		// Échap ne démarre pas une partie depuis l'accueil ; il bascule la pause.
		echapAppuye := ordta.KeyState().Pressed(window.KeyEscape)
		if echapAppuye && !echapEtaitAppuye && menu.Demarre {
			if menu.Ouvert {
				reprendreJeu()
			} else {
				menu.Pause()
				sons.ArreterEffets()
				interfaceJeu.SetVisible(false)
				gui.Manager().Set(menu.Panneau)
				LibererSourisSimulation()
			}
		}
		echapEtaitAppuye = echapAppuye
		ambiance := "monde"
		if menu.Ouvert {
			ambiance = "menu"
		} else if combatEnCours && combatBackend != nil {
			ambiance = "combat"
			if combatBackend.Mode == library.ModeArene && combatBackend.NumeroVague == len(combatBackend.Vagues) {
				ambiance = "boss"
			}
			if combatBackend.Phase == library.PhaseVictoire {
				ambiance = "victoire"
			}
			if combatBackend.Phase == library.PhaseDefaite {
				ambiance = "defaite"
			}
		}
		sons.Ambiance(ambiance, float32(tempsImage.Seconds()))
		if menu.Ouvert {
			// Rien ne progresse : déplacement, animations, effets et tours de combat.
			// Mémoriser les touches empêche TAB/E/F5 de se déclencher à la reprise.
			tabEtaitAppuye = ordta.KeyState().Pressed(window.KeyTab)
			eEtaitAppuye = ordta.KeyState().Pressed(window.KeyE)
			f5EtaitAppuye = ordta.KeyState().Pressed(window.KeyF5)
			ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
			rendu.Render(scene, cameraSimulation)
			return
		}
		deplacementEffectue := false
		for _, resultat := range effetsPersonnage.MettreAJour(tempsImage) {
			interfaceInventaire.AfficherMessage(resultat.Message)
		}
		personnage3d.Equiper(PersonnageBackend.Equipement.Tete.Nom != "", PersonnageBackend.Equipement.Torse.Nom != "", PersonnageBackend.Equipement.Pied.Nom != "")
		tempsSauvegarde += tempsImage
		f5Appuye := ordta.KeyState().Pressed(window.KeyF5)
		if (!combatEnCours && tempsSauvegarde >= 15*time.Second) || (f5Appuye && !f5EtaitAppuye) {
			sauvegarder()
			tempsSauvegarde = 0
		}
		f5EtaitAppuye = f5Appuye

		// TAB fonctionne comme un interrupteur, une seule fois par pression.
		tabAppuye := ordta.KeyState().Pressed(window.KeyTab)
		if tabAppuye && !tabEtaitAppuye {
			if interfaceInventaire.Ouvert {
				fermerInventaire()
			} else {
				ouvrirInventaire()
			}
		}
		tabEtaitAppuye = tabAppuye
		eAppuye := ordta.KeyState().Pressed(window.KeyE)
		interactionDemandee := eAppuye && !eEtaitAppuye
		eEtaitAppuye = eAppuye

		// 1. Déplacement autorisé uniquement hors course.
		if !combatEnCours && !interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert {
			positionActuelle := noeudPersonnage.Position()
			vitessePersonnage := vitesseMarche

			// Distance parcourue pendant cette frame.
			distance := vitessePersonnage * float32(tempsImage.Seconds())

			// Direction demandée par le joueur.
			commandeAvant, commandeDroite := float32(0), float32(0)

			if ordta.KeyState().Pressed(window.KeyW) {
				commandeAvant += 1
			}

			if ordta.KeyState().Pressed(window.KeyS) {
				commandeAvant -= 1
			}

			if ordta.KeyState().Pressed(window.KeyA) {
				commandeDroite -= 1
			}

			if ordta.KeyState().Pressed(window.KeyD) {
				commandeDroite += 1
			}
			if ordta.KeyState().Pressed(window.KeyLeftShift) {
				vitessePersonnage = vitesseSprint
				distance = vitessePersonnage * float32(tempsImage.Seconds())
			}

			directionX, directionY := CalculerDirectionPersonnage(angleHorizontal, commandeAvant, commandeDroite)

			// Calcul de la nouvelle position dans le repère du monde.
			positionFinaleX, positionFinaleY := deplacerAvecCollisions(positionActuelle.X, positionActuelle.Y, directionX*distance, directionY*distance, rayonPersonnage, donneesMonde)

			noeudPersonnage.SetPosition(positionFinaleX, positionFinaleY, positionActuelle.Z)
			deplacementEffectue = positionFinaleX != positionActuelle.X || positionFinaleY != positionActuelle.Y

			// 2. Vérification après le déplacement.
			if EstDansArene(positionFinaleX, positionFinaleY) {
				combatEnCours = true

				noeudPersonnage.SetPosition(
					PositionJoueurArene.X,
					PositionJoueurArene.Y,
					PositionJoueurArene.Z,
				)

				noeudPersonnage.SetRotationZ(math32.Pi / 2)

				LibererSourisSimulation()

				// L'entrée ouvre le choix : aucun monstre n'attaque avant validation.
				choixCombat.Ouvrir()

			}
		}

		// Un PNJ ne peut être utilisé qu'en exploration et à courte distance.
		var pnjProche *PNJ
		if !combatEnCours && !interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert {
			pnjProche = PNJLePlusProche(noeudPersonnage.Position(), pnjs)
		}

		if interactionDemandee {
			if interfaceCommerce.Ouvert {
				fermerCommerce()
			} else if pnjProche != nil {
				pnjActif = pnjProche
				interfaceCommerce.Ouvrir(pnjProche.Type)
				LibererSourisSimulation()
			}
		}

		if interfaceCommerce.Ouvert || combatEnCours || interfaceInventaire.Ouvert {
			interfaceInteraction.Afficher(nil)
		} else {
			interfaceInteraction.Afficher(pnjProche)
		}
		interfaceCommerce.MettreAJour()

		// 3. La caméra est choisie après la détection.
		if interfaceInventaire.Ouvert {
			PlacerCameraInventaire(cameraSimulation, noeudPersonnage)
			noeudPersonnage.SetRotationZ(0)
		} else if interfaceCommerce.Ouvert && pnjActif != nil {
			PlacerCameraDialogue(cameraSimulation, noeudPersonnage, pnjActif.Position)
		} else if combatEnCours {
			PlacerCameraArene(cameraSimulation)
			noeudPersonnage.SetRotationZ(math32.Pi / 2)
		} else {
			// 4. Caméra personnage
			PlacerCameraPersonnage(cameraSimulation, noeudPersonnage, angleHorizontal, angleVertical)

			// En exploration seulement, le personnage suit le regard de la souris.
			angle := math32.Pi - angleHorizontal
			noeudPersonnage.SetRotationZ(angle)
		}

		// Mise à jour du combat à chaque image.
		if combatEnCours &&
			combatBackend != nil &&
			combatBackend.Phase == library.PhaseTourMonstres && !combatBackend.EffetEnCours3D() &&
			!interfaceInventaire.Ouvert && !animationCombat.Active {

			tempsAvantActionMonstre += tempsImage

			if tempsAvantActionMonstre >= time.Second {
				indexActeur := combatBackend.IndexMonstreActif
				for indexActeur < len(combatBackend.Ennemis) && combatBackend.Ennemis[indexActeur].Monstre.PvActuel <= 0 {
					indexActeur++
				}
				defenseAvant := combatBackend.DefenseActive
				tourAvant := combatBackend.Tour
				resultat := combatBackend.ProchaineActionMonstre()
				if resultat.Type == "attaque_monstre" {
					if defenseAvant {
						resultat.Message += " Protection consommée : dégâts réduits de moitié."
					}
					if indexActeur < len(combatBackend.Ennemis) {
						genre := combatBackend.Ennemis[indexActeur].Type
						if tourAvant%3 == 0 && (genre == library.TypeGobelin || genre == library.TypeTroll) {
							resultat.Message += " Coup renforcé : puissance doublée."
						}
					}
				}
				journalCombat.Ajouter(texteResultatCombat("Action monstre", resultat))
				if resultat.Reussite && indexActeur < len(modelesAdversaires) {
					acteur := modelesAdversaires[indexActeur]
					cible := noeudPersonnage
					action := string(combatBackend.Ennemis[indexActeur].Type)
					if resultat.Type == "soin_monstre" {
						action = "soin_monstre"
						cible = acteur.Noeud()
						for index, ennemi := range combatBackend.Ennemis {
							if ennemi.Monstre.Nom == resultat.Cible {
								cible = modelesAdversaires[index].Noeud()
							}
						}
					}
					animationCombat.Demarrer(action, acteur.Noeud(), cible, acteur.Membres, nil)
				}
				if resultat.Degats > 0 {
					sons.Effet("impact")
				}
				dernierMessageCombat = resultat.Message

				fmt.Println(resultat.Message)
				fmt.Println(
					"Vie :",
					PersonnageBackend.PvActuel,
					"/",
					PersonnageBackend.PvMaxTotal,
				)

				tempsAvantActionMonstre = 0
			}
		} else {
			tempsAvantActionMonstre = 0
		}

		// Effets temporaires non bloquants, par exemple la potion de poison.
		if combatEnCours && combatBackend != nil && !interfaceInventaire.Ouvert {
			for _, resultat := range combatBackend.MettreAJourEffets(tempsImage) {
				dernierMessageCombat = resultat.Message
				journalCombat.Ajouter(texteResultatCombat("Effet temporaire", resultat))
			}
		}

		// Une courte transition sépare deux vagues.
		if combatEnCours &&
			combatBackend != nil &&
			combatBackend.Phase == library.PhaseEntreVagues && !animationCombat.Active {

			tempsAvantVagueSuivante += tempsImage

			if tempsAvantVagueSuivante >= 2*time.Second {
				for _, modele := range modelesAdversaires {
					scene.Remove(modele.Noeud())
				}

				if err := combatBackend.CommencerVagueSuivante(); err != nil {
					panic(err)
				}

				// L’initiative choisie par le backend est conservée.
				combatBackend.IndexMonstreActif = 0
				cibleSelectionnee = 0
				interfaceCombat.MasquerChoix()

				nouveauxModeles, err := ChargerMontres(scene, combatBackend.Ennemis)
				if err != nil {
					panic(err)
				}
				modelesAdversaires = nouveauxModeles

				dernierMessageCombat = fmt.Sprintf(
					"Vague %d : début du combat.",
					combatBackend.NumeroVague,
				)
				tempsAvantVagueSuivante = 0
			}
		} else {
			tempsAvantVagueSuivante = 0
		}

		// Visibilité de l’interface
		afficherCommandes := combatEnCours &&
			combatBackend != nil &&
			combatBackend.Phase == library.PhaseTourJoueur &&
			!interfaceInventaire.Ouvert && !animationCombat.Active
		combatTermine := combatEnCours && combatBackend != nil &&
			(combatBackend.Phase == library.PhaseVictoire || combatBackend.Phase == library.PhaseDefaite) &&
			!interfaceInventaire.Ouvert && !animationCombat.Active

		interfaceCombat.Panneau.SetVisible(combatEnCours && combatBackend != nil && !interfaceInventaire.Ouvert)
		interfaceCombat.BoutonAttaquer.SetVisible(afficherCommandes)
		interfaceCombat.BoutonSorts.SetVisible(afficherCommandes)
		interfaceCombat.BoutonInventaire.SetVisible(afficherCommandes)
		interfaceCombat.BoutonDefendre.SetVisible(afficherCommandes)
		interfaceCombat.BoutonQuitter.SetVisible(combatEnCours && combatBackend != nil)
		interfaceCombat.BoutonRejouer.SetVisible(combatTermine)

		// Vie et Statistiques
		interfacePersonnage.MettreAJourPersonnage(&PersonnageBackend)
		interfacePersonnage.Panneau.SetVisible(!interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert)
		interfaceInventaire.MettreAJour(&PersonnageBackend)

		// Affichage panel des monstres quand le personnage est en combat
		if combatEnCours && combatBackend != nil && !interfaceInventaire.Ouvert {
			interfaceMonstres.MettreAJourMonstres(
				combatBackend.Ennemis,
			)
		} else {
			interfaceMonstres.MettreAJourMonstres(nil)
		}

		// État général du combat : vague, tour, phase et dernier événement.
		if combatEnCours && combatBackend != nil && !interfaceInventaire.Ouvert {
			interfaceEtatCombat.MettreAJour(
				combatBackend,
				dernierMessageCombat,
			)
		} else {
			interfaceEtatCombat.MettreAJour(nil, "")
		}

		// 6. Animations du personnage.
		if animationCombat.Active {
			animationCombat.MettreAJour(float32(tempsImage.Seconds()))
		} else {
			personnage3d.Animer(float32(tempsImage.Seconds()), deplacementEffectue)
		}
		journalCombat.Panneau.SetVisible(combatEnCours && combatBackend != nil && !interfaceInventaire.Ouvert)
		if combatBackend != nil && !animationCombat.Active {
			journalCombat.Etat(combatBackend)
		}

		// 7. Effacer l'image précédente, puis dessiner la nouvelle scène.
		ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
		rendu.Render(scene, cameraSimulation)
	})
}
