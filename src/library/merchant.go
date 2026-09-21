package library

import (
	"fmt"
)

// merchant affiche les objets disponibles chez le marchand.
func (c *Character) merchant() {
	for {
		fmt.Println("\n=== MARCHAND ===")
		fmt.Printf("1. %s - Gratuit\n", ItemPotionDeVie)
		fmt.Printf("2. %s - Gratuit\n", ItemPotionDePoison)
		fmt.Printf("3. %s - Gratuit\n", ItemLivreBouleDeFeu)
		fmt.Println("0. Retour")

		choice, ok := ReadChoice("Votre choix : ")

		if !ok {
			fmt.Println("Choix invalide. Entrez 0, 1, 2 ou 3.")
			continue
		}

		switch choice {
		case 1:
			c.buy(ItemPotionDeVie)
		case 2:
			c.buy(ItemPotionDePoison)
		case 3:
			c.buy(ItemLivreBouleDeFeu)
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
	if !c.AddInventory(item) {
		fmt.Printf("Inventaire plein (%d / %d) : impossible d'ajouter %s.\n",
			c.TotalInventaire(), c.CapaciteInventaire, item)
		return
	}
	fmt.Println("Vous avez acheté :", item)
}
