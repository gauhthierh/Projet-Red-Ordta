package main

import (
	"fmt"
	"time"
)

// Noms des objets : définis une seule fois pour éviter les fautes de frappe
// dans les clés de l'inventaire.
const (
	itemPotionDeVie = "Potion de vie"
)

func (c *Character) takePot() {
	potionDeVie := c.Inventaire[itemPotionDeVie]
	if potionDeVie > 0 {
		c.Inventaire[itemPotionDeVie]--
		c.PVActuel += 50
		if c.PVActuel >= c.PVMax {
			c.PVActuel = c.PVMax
		}
		fmt.Printf("PV : %d / %d\n", c.PVActuel, c.PVMax)
	} else {
		fmt.Println("Aucunes potions dans l'inventaire")
	}
}

func (c *Character) poisonPot() {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMax)
		c.isDead()
	}
}
