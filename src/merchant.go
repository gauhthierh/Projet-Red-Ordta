package main

import "fmt"

// merchant affiche les objets disponibles chez le marchand.
func (c *Character) merchant() {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Printf("1. %s - Gratuit\n", itemPotionDeVie)
		fmt.Println("0. Retour")

		choice, ok := readChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide. Entrez 0 ou 1.")
			continue
		}

		switch choice {
		case 1:
			item := itemPotionDeVie
			c.addInventory(item)
			fmt.Println("Vous avez acheté :", item)
		case 0:
			return
		default:
			fmt.Println("Choix invalide. Entrez 0 ou 1.")
		}
	}
}
