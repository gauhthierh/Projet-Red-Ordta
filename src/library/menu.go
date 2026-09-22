package library

import (
	"fmt"
)

/* Ce fichier gère le menu principal du jeu et le retour depuis les différentes fonctionnalités. */

/* La méthode MainMenu affiche les fonctionnalités disponibles et dirige le joueur vers celle qu'il choisit. */
func (c *Character) MainMenu() {
	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Forgeron")
		fmt.Println("5. Entrainement")
		fmt.Println("6. Qui sont ils ?")
		fmt.Println("0. Quitter")

		choice, ok := ReadChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide !")
			continue
		}

		switch choice {
		case 1:
			c.displayInfo()
			WaitForReturn()
		case 2:
			c.InventoryMenu()
		case 3:
			c.merchant()
		case 4:
			c.forgeron()
		case 5:
			c.trainingFight()
		case 6:
			fmt.Println("ABBA et Steven Spielberg")
			WaitForReturn()
		case 0:
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide !")
		}
	}
}

/* La fonction WaitForReturn attend que le joueur choisisse zéro avant de revenir au menu principal. */
func WaitForReturn() {
	for {
		fmt.Println("\n0. Retour")

		choice, ok := ReadChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide. Entrez 0 pour revenir au menu principal.")
			continue
		}

		if choice == 0 {
			return
		}

		fmt.Println("Choix invalide. Entrez 0 pour revenir au menu principal.")
	}
}
