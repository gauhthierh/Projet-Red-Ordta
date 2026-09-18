package main

import "fmt"

// addInventory ajoute un exemplaire d'un objet dans l'inventaire.
func addInventory(character *Character, item string) {
	character.Inventaire[item]++
}

// merchant affiche les objets disponibles chez le marchand.
func merchant(character *Character) {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Println("1. Potion de vie - Gratuit")
		fmt.Println("0. Retour")
		fmt.Print("Votre choix : ")

		var choice int
		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("Choix invalide. Entrez 0 ou 1.")
			clearInput()
			continue
		}

		switch choice {
		case 1:
			item := "potion de vie"
			addInventory(character, item)
			fmt.Println("Vous avez acheté :", item)
		case 0:
			return
		default:
			fmt.Println("Choix invalide. Entrez 0 ou 1.")
		}
	}
}
