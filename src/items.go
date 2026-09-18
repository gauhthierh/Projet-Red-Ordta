package main

import "fmt"

func (c *Character) takePot() {
	potionDeVie := c.Inventaire["potion de vie"]
	if potionDeVie > 0 {
		c.Inventaire["potion de vie"]--
		c.PVActuel += 50
		if c.PVActuel >= c.PVMax {
			c.PVActuel = c.PVMax
		}
		fmt.Printf("PV : %d / %d\n", c.PVActuel, c.PVMax)
	} else {
		fmt.Println("Aucunes potions dans l'inventaire")
	}
}
