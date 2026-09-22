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
	partieChargee, err := ChargerPartie3D(cheminSauvegarde)
	if err != nil {
		panic(err)
	}
	if partieChargee != nil {
		PersonnageBackend = partieChargee.Personnage
		position := partieChargee.Position
		if positionAutorisee(position.X, position.Y, rayonPersonnage, donneesMonde) && !EstDansArene(position.X, position.Y) {
			spawn = [3]float32{position.X, position.Y, position.Z}
		}
	}

	var combatBackend *library.CombatArene
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
	interfaceCombat := NouvelleInterfaceCombat(scene)
	interfacePersonnage := NouvelleInterfacePersonnage(scene)
	interfaceMonstres := NouvelleInterfaceMonstres(scene, len(PositionsMonstresArene))
	interfaceEtatCombat := NouvelleInterfaceEtatCombat(scene)
	interfaceInventaire := NouvelleInterfaceInventaire(scene)
	interfaceInteraction := NouvelleInterfaceInteraction(scene)
	interfaceCommerce := NouvelleInterfaceCommerce(scene, &PersonnageBackend)
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

	VerrouillerSourisSimulation()
	angleHorizontal := float32(0)
	angleVertical := float32(0)
	if partieChargee != nil {
		angleHorizontal = partieChargee.AngleHorizontal
		angleVertical = partieChargee.AngleVertical
	}
	tempsSauvegarde := time.Duration(0)
	f5EtaitAppuye := false
	sauvegarder := func() {
		if combatEnCours {
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
	defer sauvegarder()

	ActiverRegardSouris(ordta, &angleHorizontal, &angleVertical, func() bool {
		return !combatEnCours && !interfaceInventaire.Ouvert && !interfaceCommerce.Ouvert
	})

	fermerInventaire := func() {
		interfaceInventaire.Fermer()
		sauvegarder()
		if !combatEnCours {
			VerrouillerSourisSimulation()
		}
	}

	ouvrirInventaire := func() {
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

	appliquerResultat := func(resultat library.ResultatAction, indexCible int) {
		dernierMessageCombat = resultat.Message
		if resultat.CibleVaincue &&
			indexCible >= 0 &&
			indexCible < len(modelesAdversaires) {
			modelesAdversaires[indexCible].Noeud().SetVisible(false)
		}
		if resultat.Reussite && resultat.TourConsomme {
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
			if combatBackend.Ennemis[indexMonstre].Monstre.PVActuel <= 0 {
				return
			}
			cibleSelectionnee = indexMonstre
			dernierMessageCombat = "Cible : " + combatBackend.Ennemis[indexMonstre].Monstre.Nom
		})
	}

	interfaceCombat.BoutonAttaquer.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if combatBackend != nil && combatBackend.Phase == library.PhaseTourJoueur {
			interfaceCombat.AfficherChoix("attaque", combatBackend.AttaquesDisponibles())
		}
	})

	interfaceCombat.BoutonSorts.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if combatBackend != nil && combatBackend.Phase == library.PhaseTourJoueur {
			interfaceCombat.AfficherChoix("sort", combatBackend.SortsDisponibles())
		}
	})

	interfaceCombat.BoutonInventaire.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		ouvrirInventaire()
	})

	interfaceCombat.BoutonDefendre.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
		if combatBackend == nil {
			return
		}
		resultat := combatBackend.Defendre()
		appliquerResultat(resultat, -1)
	})

	for index, bouton := range interfaceCombat.BoutonsChoix {
		indexOption := index
		bouton.Subscribe(gui.OnClick, func(nomEvenement string, evenement interface{}) {
			if combatBackend == nil ||
				combatBackend.Phase != library.PhaseTourJoueur ||
				indexOption >= len(interfaceCombat.OptionsChoix) {
				return
			}

			cibles := combatBackend.CiblesVivantes()
			if len(cibles) == 0 {
				return
			}
			if cibleSelectionnee >= len(combatBackend.Ennemis) ||
				combatBackend.Ennemis[cibleSelectionnee].Monstre.PVActuel <= 0 {
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
			appliquerResultat(resultat, cibleSelectionnee)
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
			if combatBackend == nil {
				return
			}
			resultat := combatBackend.UtiliserObjet(nomObjet)
			interfaceInventaire.AfficherMessage(resultat.Message)
			appliquerResultat(resultat, -1)
			if resultat.Reussite && resultat.TourConsomme {
				fermerInventaire()
			}
			return
		}

		resultat := PersonnageBackend.UtiliserObjet3D(nomObjet)
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

	// Le bouton de sortie n'est utilisable qu'après une victoire ou une défaite.
	interfaceCombat.BoutonQuitter.Subscribe(
		gui.OnClick,
		func(nomEvenement string, evenement interface{}) {
			if combatBackend == nil {
				return
			}

			combatTermine := combatBackend.Phase == library.PhaseVictoire ||
				combatBackend.Phase == library.PhaseDefaite
			if !combatTermine {
				return
			}

			// Après une défaite, le joueur ressort avec la moitié de sa vie.
			if combatBackend.Phase == library.PhaseDefaite {
				combatBackend.RessusciterApresDefaite()
			}

			for _, modele := range modelesAdversaires {
				scene.Remove(modele.Noeud())
			}

			modelesAdversaires = nil
			combatBackend = nil
			combatEnCours = false
			cibleSelectionnee = 0
			dernierMessageCombat = ""
			tempsAvantActionMonstre = 0
			tempsAvantVagueSuivante = 0

			noeudPersonnage.SetPosition(
				PositionSortieArene.X,
				PositionSortieArene.Y,
				PositionSortieArene.Z,
			)

			interfaceCombat.MasquerChoix()
			interfaceInventaire.Fermer()
			sauvegarder()
			VerrouillerSourisSimulation()
		},
	)

	// Lumière dans le jeu
	scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 1))

	// Boucle principale : mise à jour du jeu, puis affichage de chaque image.
	ordta.Gls().ClearColor(0.15, 0.25, 0.30, 1)
	ordta.Run(func(rendu *renderer.Renderer, tempsImage time.Duration) {
		deplacementEffectue := false
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

				// Système de combat
				nouveuCombat, err := library.NouveauCombatArene(&PersonnageBackend)

				if err != nil {
					panic(err)
				}

				combatBackend = nouveuCombat

				// Le joueur commence toujours la première vague.
				combatBackend.Phase = library.PhaseTourJoueur
				combatBackend.IndexMonstreActif = 0
				cibleSelectionnee = 0
				interfaceCombat.MasquerChoix()
				dernierMessageCombat = "Le combat commence. À vous de jouer."

				modelesVague, err := ChargerMontres(scene, combatBackend.Ennemis)
				if err != nil {
					panic(err)
				}
				modelesAdversaires = modelesVague

				fmt.Println("Vague :", combatBackend.NumeroVague)
				fmt.Println("Phase :", combatBackend.Phase)
				fmt.Println("Modèles 3D chargés :", len(modelesAdversaires))

				for _, adversaire := range combatBackend.Ennemis {
					fmt.Println(
						adversaire.Monstre.Nom,
						adversaire.Monstre.PVActuel,
						"/",
						adversaire.Monstre.PVMax,
					)
				}

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
			combatBackend.Phase == library.PhaseTourMonstres &&
			!interfaceInventaire.Ouvert {

			tempsAvantActionMonstre += tempsImage

			if tempsAvantActionMonstre >= time.Second {
				resultat := combatBackend.ProchaineActionMonstre()
				dernierMessageCombat = resultat.Message

				fmt.Println(resultat.Message)
				fmt.Println(
					"Vie :",
					PersonnageBackend.PVActuel,
					"/",
					PersonnageBackend.PVMaxTotal,
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
			}
		}

		// Une courte transition sépare deux vagues.
		if combatEnCours &&
			combatBackend != nil &&
			combatBackend.Phase == library.PhaseEntreVagues {

			tempsAvantVagueSuivante += tempsImage

			if tempsAvantVagueSuivante >= 2*time.Second {
				for _, modele := range modelesAdversaires {
					scene.Remove(modele.Noeud())
				}

				if err := combatBackend.CommencerVagueSuivante(); err != nil {
					panic(err)
				}

				// Comme à l'entrée de l'arène, le joueur commence la vague.
				combatBackend.Phase = library.PhaseTourJoueur
				combatBackend.IndexMonstreActif = 0
				cibleSelectionnee = 0
				interfaceCombat.MasquerChoix()

				nouveauxModeles, err := ChargerMontres(scene, combatBackend.Ennemis)
				if err != nil {
					panic(err)
				}
				modelesAdversaires = nouveauxModeles

				dernierMessageCombat = fmt.Sprintf(
					"Vague %d : à vous de jouer.",
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
			!interfaceInventaire.Ouvert
		combatTermine := combatEnCours && combatBackend != nil &&
			(combatBackend.Phase == library.PhaseVictoire || combatBackend.Phase == library.PhaseDefaite) &&
			!interfaceInventaire.Ouvert

		interfaceCombat.Panneau.SetVisible(afficherCommandes || combatTermine)
		interfaceCombat.BoutonAttaquer.SetVisible(afficherCommandes)
		interfaceCombat.BoutonSorts.SetVisible(afficherCommandes)
		interfaceCombat.BoutonInventaire.SetVisible(afficherCommandes)
		interfaceCombat.BoutonDefendre.SetVisible(afficherCommandes)
		interfaceCombat.BoutonQuitter.SetVisible(combatTermine)

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
		personnage3d.Animer(float32(tempsImage.Seconds()), deplacementEffectue)

		// 7. Effacer l'image précédente, puis dessiner la nouvelle scène.
		ordta.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
		rendu.Render(scene, cameraSimulation)
	})
}
