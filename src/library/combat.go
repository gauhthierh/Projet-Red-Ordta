package library

import (
	"fmt"
	"math/rand"
)

// ================================================= //
// === Ce fichier gère le déroulement du combat  === //
// ================================================= //

// Lance un combat d'entraînement contre un gobelin //
// Le combat est au tour par tour //
func (c *Character) TrainingFight() {
	adversaire := InitGoblin(c.Niveau)
	tour := 1
	c.Initiative = rand.Intn(10) + 1
	adversaire.Initiative = rand.Intn(10) + 1
	fmt.Printf("\nDébut du combat contre %s (niveau %d) !\n", adversaire.Nom, adversaire.Niveau)
	fmt.Printf("Initiative de %s : %d\n", c.Nom, c.Initiative)
	fmt.Printf("Initiative de %s : %d\n", adversaire.Nom, adversaire.Initiative)
	personnageCommence := c.Initiative >= adversaire.Initiative
	if personnageCommence {
		fmt.Printf("%s commence le combat !\n", c.Nom)
	} else {
		fmt.Printf("%s commence le combat !\n", adversaire.Nom)
	}
	for {
		fmt.Printf("\n=== TOUR %d ===\n", tour)
		if personnageCommence {
			c.CharacterTurn(&adversaire)
			if c.FinCombat(&adversaire) {
				return
			}
			adversaire.GoblinPattern(c, tour)
			if c.FinCombat(&adversaire) {
				return
			}
		} else {
			adversaire.GoblinPattern(c, tour)
			if c.FinCombat(&adversaire) {
				return
			}
			c.CharacterTurn(&adversaire)
			if c.FinCombat(&adversaire) {
				return
			}
		}
		tour++
	}
}

// Affiche le menu de combat du joueur //
// Persiste jusqu'a ce qu'il ait joué //
func (c *Character) CharacterTurn(m *Monster) {
	for {
		fmt.Println("=== COMBAT ===")
		fmt.Println("1. Attaque physique")
		fmt.Println("2. Sort")
		fmt.Println("3. Inventaire")
		choix := ReadChoiceEntre("Votre choix : ", 3)
		switch choix {
		case 1:
			if c.ChoixAttaquePhysique(m) {
				return
			}
		case 2:
			if c.ChoixSort(m) {
				return
			}
		case 3:
			if c.ChoixInventaire() {
				return
			}
		default:
			fmt.Println("Choix invalide, veuillez entrer un choix valide")
		}
	}
}

// Vérifie si le joueur ou le monstre est à 0 Pv //
// Annonce le résultat et renvoie true si le combat est terminé //
func (c *Character) FinCombat(m *Monster) bool {
	if m.PvActuel <= 0 {
		fmt.Printf("%s est vaincu !\n", m.Nom)
		fmt.Println("Vous avez gagné l'entraînement, bien joué !")
		gainOr := m.OrMin + rand.Intn(m.OrMax-m.OrMin+1)
		c.Argent += gainOr
		fmt.Printf("Vous trouvez %d pièces d'or. Bourse : %d Po\n", gainOr, c.Argent)
		c.GagnerExperience(m.ExperienceDonnee)
		return true
	}
	if c.PvActuel <= 0 {
		c.IsDead()
		fmt.Println("Le combat d'entraînement est terminé.")
		return true
	}
	return false
}

// Affiche l'inventaire en combat et utilise l'objet choisi //
// L'équipement ne peut être choisi en combat //
// Renvoie true si l'objet a eu un effet //
func (c *Character) ChoixInventaire() bool {
	for {
		objets := c.SortedItems()
		if len(objets) == 0 {
			fmt.Println("L'inventaire est vide !")
			return false
		}
		for indice, objet := range objets {
			fmt.Printf("%d. %s, quantité : %d\n", indice+1, objet, c.Inventaire[objet])
		}
		fmt.Println("0. Retour")
		choix := ReadChoiceEntre("Votre choix : ", len(objets))
		if choix == 0 {
			return false
		}
		objet := objets[choix-1]
		if _, estEquipement := TrouverEquipement(objet); estEquipement {
			fmt.Println("L'équipement ne peut pas être changé pendant le combat.")
			continue
		}
		fmt.Printf("Vous utilisez %s\n", objet)
		if c.UseItem(objets[choix-1]) {
			return true
		} else {
			fmt.Println("Choisissez un autre objet ou 0 pour revenir en arrière")
		}
	}
}
