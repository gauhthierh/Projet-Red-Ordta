package main

import (
	"fmt"
)

func (c *Character) isDead() {
	if c.PVActuel <= 0 {
		fmt.Println(c.Nom, "est mort!")
		c.PVActuel = c.PVMax / 2
		fmt.Println(c.Nom, "a maintenant", c.PVActuel, "points de vie")
	} else {
		fmt.Println(c.Nom, "est vivant")
	}
}
