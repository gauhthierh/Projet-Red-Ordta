package library

import (
	"fmt"
)

/*
mainMenu affiche les choix principaux du jeu.
Le Pointeur permet aux autres menus de modifier le personnage
*/
func (c *Character) MainMenu() {
	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Forgeron")
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
		case 0:
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide !")
		}
	}
}

// waitForReturn laisse le temps de lire l'écran avant le retour au menu.
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
