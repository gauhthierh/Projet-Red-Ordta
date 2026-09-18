package main

import "fmt"

func (c Character) accessInventory() {
	fmt.Println("Inventaire :")
	vide := true
	for nom, quantite := range c.Inventaire {
		if quantite > 0 {
			fmt.Printf("- %s : %d\n", nom, quantite)
			vide = false
		}
	}
	if vide {
		fmt.Println("(inventaire vide)")
	}
}

// addInventory ajoute un exemplaire d'un objet dans l'inventaire.
func (c *Character) addInventory(item string) {
	c.Inventaire[item]++
}
