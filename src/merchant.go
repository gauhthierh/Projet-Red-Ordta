package main

import "fmt"

// merchant affiche les objets disponibles chez le marchand.
func (c *Character) merchant() {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Printf("1. %s - Gratuit\n", itemPotionDeVie)
		fmt.Printf("2. %s - Gratuit\n", itemPotionDePoison)
		fmt.Printf("3. %s - Gratuit\n", itemLivreBouleDeFeu)
		fmt.Println("0. Retour")

		choice, ok := readChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide. Entrez 0, 1, 2 ou 3.")
			continue
		}

		switch choice {
		case 1:
			item := itemPotionDeVie
			c.addInventory(item)
			fmt.Println("Vous avez acheté :", item)
		case 2:
			item := itemPotionDePoison
			c.addInventory(item)
			fmt.Println("Vous avez acheté :", item)
		case 3:
			item := itemLivreBouleDeFeu
			c.addInventory(item)
			fmt.Println("Vous avez acheté :", item)
		case 0:
			return
		default:
			fmt.Println("Choix invalide. Entrez 0, 1, 2 ou 3.")
		}
	}
}
