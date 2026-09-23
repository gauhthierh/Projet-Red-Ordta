package library

import (
	"fmt"
	"slices"
)

// ============================================== //
// === Ce fichier gère les attaques physiques === //
// ============================================== //

// Affiche les attaques physiques connues et attaque le monstre avec //
// renvoie true si une attaque a été lancée //
func (c *Character) ChoixAttaquePhysique(m *Monster) bool {
	if len(c.AttaquesPhysiques) == 0 {
		fmt.Println("Vous ne connaissez aucune attaque physique.")
		return false
	}
	for {
		fmt.Println("\n=== ATTAQUES PHYSIQUES ===")
		for indice, attaque := range c.AttaquesPhysiques {
			degats, connu := c.CalculerDegatsPhysique(attaque)
			if !connu {
				fmt.Printf("%d. %s (indisponible)\n", indice+1, attaque)
			} else {
				fmt.Printf("%d. %s (%d dégâts)\n", indice+1, attaque, degats)
			}
		}
		fmt.Println("0. Retour")
		choix := ReadChoiceEntre("Votre choix : ", len(c.AttaquesPhysiques))
		if choix == 0 {
			return false
		}
		attaqueChoisi := c.AttaquesPhysiques[choix-1]
		degats, connu := c.CalculerDegatsPhysique(attaqueChoisi)
		if !connu {
			fmt.Println("Cette attaque ne peut pas être utilisée.")
			continue
		}
		m.SubirDegats(degats)
		fmt.Printf("%s utilise %s sur %s et lui inflige %d dégâts !\n", c.Nom, attaqueChoisi, m.Nom, degats)
		fmt.Printf("Pv de %s : %d / %d\n", m.Nom, m.PvActuel, m.PvMax)
		return true
	}
}

// Renvoie les dégâts de base d'une attaque physique //
// false si l'attaque est inconnue //
func InfosAttaquePhysique(attaque string) (int, bool) {
	switch attaque {
	case AttaqueBasique:
		return DegatsAttaqueBasique, true
	case AttaqueClaquounette:
		return DegatsClaquounette, true
	case AttaqueCoupsDePied:
		return DegatsCoupsDePied, true
	case AttaquePichenette:
		return DegatsPichenette, true
	case AttaqueMorsure:
		return DegatsMorsure, true
	case AttaqueUppercut:
		return DegatsUppercut, true
	default:
		return 0, false
	}
}

// Renvoie les dégâts d'une attaque physique //
// et ajoute le bonus d'attaque du joueur //
func (c *Character) CalculerDegatsPhysique(attaque string) (int, bool) {
	degats, existe := InfosAttaquePhysique(attaque)
	if !existe {
		return 0, false
	} else {
		return degats + c.Attaque, existe
	}
}

// Apprend une attaque physique au joueur //
// renvoie false s'il la connaît déjà //
func (c *Character) ApprentissageAttaque(attaque string) bool {
	if slices.Contains(c.AttaquesPhysiques, attaque) {
		return false
	}
	c.AttaquesPhysiques = append(c.AttaquesPhysiques, attaque)
	return true
}

// Renvoie l'attaque enseignée par un manuel de combat //
// false si l'objet n'est pas un manuel //
func AttaqueDuManuel(objet string) (string, bool) {
	switch objet {
	case ItemLivreAttaqueClaquounette:
		return AttaqueClaquounette, true
	case ItemLivreAttaqueCoupsDePied:
		return AttaqueCoupsDePied, true
	case ItemLivreAttaquePichenette:
		return AttaquePichenette, true
	case ItemLivreAttaqueMorsure:
		return AttaqueMorsure, true
	case ItemLivreAttaqueUppercut:
		return AttaqueUppercut, true
	default:
		return "", false
	}
}
