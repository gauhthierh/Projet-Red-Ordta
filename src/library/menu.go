package library

import (
	"fmt"
)

// ================================================ //
// === Ce fichier gère le menu principal du jeu === //
// ================================================ //

// Affiche le menu principal et dirige le joueur vers la fonctionnalité choisie //
// Le jeu s'arrête quand le joueur choisit de quitter //
func (c *Character) MainMenu() {
	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Forgeron")
		fmt.Println("5. Entraînement")
		fmt.Println("6. Qui sont-ils ?")
		fmt.Println("0. Quitter")
		choice := ReadChoiceEntre("Votre choix : ", 6)
		switch choice {
		case 1:
			c.displayInfo()
			WaitForReturn()
		case 2:
			c.InventoryMenu()
		case 3:
			c.Merchant()
		case 4:
			c.Forgeron()
		case 5:
			c.TrainingFight()
		case 6:
			fmt.Println("ABBA et Steven Spielberg")
			WaitForReturn()
		case 0:
			fmt.Println("À bientôt !")
			return
		}
	}
}

// Attend que le joueur tape 0 pour revenir au menu, le temps de lire l'écran //
func WaitForReturn() {
	fmt.Println("\n0. Retour")
	ReadChoiceEntre("Votre choix : ", 0)
}
