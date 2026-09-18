package main

import (
	"fmt"
	"time"
)

func (c *Character) poisonPot() {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.PVActuel -= 10
		fmt.Printf("%s a été empoisonné ! PV : %d / %d\n", c.Nom, c.PVActuel, c.PVMax)
		c.isDead()
	}
}
