package main

import "fmt"

/*
mainMenu affiche les choix principaux du jeu.
Le Pointeur permet aux autres menus de modifier le personnage
*/
func (c *Character) mainMenu() {
	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Quitter")

		choice, ok := readChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide. Entrez 1, 2, 3 ou 4.")
			continue
		}

		switch choice {
		case 1:
			c.displayInfo()
			waitForReturn()
		case 2:
			c.inventoryMenu()
		case 3:
			c.merchant()
		case 4:
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide. Entrez 1, 2, 3 ou 4.")
		}
	}
}

// waitForReturn laisse le temps de lire l'écran avant le retour au menu.
func waitForReturn() {
	for {
		fmt.Println("\n0. Retour")

		choice, ok := readChoice("Votre choix : ")

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
