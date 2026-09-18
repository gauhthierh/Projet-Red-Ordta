package main

import (
	"fmt"
	"time"
)

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

func (c *Character) poisonPot() {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMax)
		c.isDead()
	}
}
