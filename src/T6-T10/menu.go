package main

import "fmt"

/*
mainMenu affiche les choix principaux du jeu.
Le Pointeur permet aux autres menus de modifier le personnage
*/
func mainMenu(character *Character) {
	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Quitter")
		fmt.Print("Votre choix : ")

		var choice int
		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("Choix invalide. Entrez 1, 2, 3 ou 4.")
			clearInput()
			continue
		}

		switch choice {
		case 1:
			character.displayInfo()
			waitForReturn()
		case 2:
			character.accessInventory()
			waitForReturn()
		case 3:
			merchant(character)
		case 4:
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide. Entrez 1, 2, 3 ou 4.")
		}
	}
}

// clearInput retire la saisie invalide avant de réafficher le menu.
func clearInput() {
	var invalidInput string
	fmt.Scanln(&invalidInput)
}

// waitForReturn laisse le temps de lire l'écran avant le retour au menu.
func waitForReturn() {
	for {
		fmt.Println("\n0. Retour")
		fmt.Print("Votre choix : ")

		var choice int
		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("Choix invalide. Entrez 0 pour revenir au menu principal.")
			clearInput()
			continue
		}

		if choice == 0 {
			return
		}

		fmt.Println("Choix invalide. Entrez 0 pour revenir au menu principal.")
	}
}
