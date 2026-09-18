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
			c.buy(itemPotionDeVie)
		case 2:
			c.buy(itemPotionDePoison)
		case 3:
			c.buy(itemLivreBouleDeFeu)
		case 0:
			return
		default:
			fmt.Println("Choix invalide. Entrez 0, 1, 2 ou 3.")
		}
	}
}

// buy ajoute l'objet acheté à l'inventaire, ou explique pourquoi c'est
// impossible : « Vous avez acheté » ne s'affiche que si l'ajout a eu lieu.
func (c *Character) buy(item string) {
	if !c.addInventory(item) {
		fmt.Printf("Inventaire plein (%d / %d) : impossible d'ajouter %s.\n",
			c.totalInventaire(), c.CapaciteInventaire, item)
		return
	}
	fmt.Println("Vous avez acheté :", item)
}
